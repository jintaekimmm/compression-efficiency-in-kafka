# Compression Efficiency in Kafka

Kafka에 데이터를 저장할 시에 압축에 따른 효율성을 확인하기 위한 시도

## Test Environment

### 1. System
```text
CPU: Intel(R) Xeon(R) CPU E3-1231 v3 @ 3.40GHz
OS : VMware base - CentOS7
Kafka : v2.7.0

 - Producer Spec : 4 Core, 8GB Memory  
 - Kafka Spec : 4 Core, 8GB Memory
```

Kafka Compression 설정은 Topic과 Producer 모두 설정
```shell
bin/kafka-topics.sh --zookeeper kafka:2181 --replication-factor 1 --partitions 1 --create --topic pb-marshaling-and-zstd --config compression.type=zstd
```
```go
func newAsyncProducer(conf *Config, compressionType string) (sarama.AsyncProducer, error) {
	... 중략
    switch compressionType {
        case "gzip":
        saramaCfg.Producer.Compression = sarama.CompressionGZIP
        case "snappy":
        saramaCfg.Producer.Compression = sarama.CompressionSnappy
        case "lz4":
        saramaCfg.Producer.Compression = sarama.CompressionLZ4
        case "zstd":
        // zstd compression requires Version >= V2_1_0_0
        saramaCfg.Version = sarama.V2_7_0_0
        saramaCfg.Producer.Compression = sarama.CompressionZSTD
        case "":
    }
}

// Producer 생성
p, err := kafka.NewProducer("zstd")
if err != nil {
    panic(err)
}
```

Kafka `linger.ms` 는 `3초`로 설정
> sarama client에서 linger.ms 설정이 아래 설정과 같은지는 확신할 수 없으나, 비슷한 동작을 하는 것으로 보인다 
```go
saramaCfg := sarama.NewConfig()
saramaCfg.Producer.Flush.Frequency = time.Second * 3
```

### 2. Sample Data
```text
CSV 파일 사용
Rows : 50,000,000
CSV Size :  6.5GB

개별 레코드는 다른 값을 가진다
```
Output Data Json Format
```json
{
    "before_mymoney": 0,
    "send_status": 1,
    "product_cnt": 1,
    "send_cnt": 1,
    "deny_cnt": 0,
    "receive_cnt": 1,
    "send_tm":"2021-10-01T22:01:01",
    "product_mymoney": 1,
    "buy_mymoney": 1,
    "discount_rate":0,
    "buy_ipaddr": "192.168.0.2",
    "buy_login_id": 12345678,
    "buy_login_name": "홍길동",
    "receive_tm": nan
}
```
### 3. Try Compression Type
Kafka에서 효율성을 확인하기 위해 여러 압축 방식을 시도하여 데이터를 저장한다

```text
 - Json(Uncompress)
 - Json + Gzip
 - Json + Lz4
 - Json + Snappy
 - Json + Zstd
```
```text
 - ProtocolBuffer
 - ProtocolBuffer + Gzip
 - ProtocolBuffer + Lz4
 - ProtocolBuffer + Snappy
 - ProtocolBuffer + Zstd
```

### 4. Test code
- [json-marshaling-only_test.go](bench/json-marshaling-only_test.go)
- [json-marshaling-gzip_test.go](bench/json-marshaling-gzip_test.go)
- [json-marshaling-lz4_test.go](bench/json-marshaling-lz4_test.go)
- [json-marshaling-snappy_test.go](bench/json-marshaling-snappy_test.go)
- [json-marshaling-zstd_test.go](bench/json-marshaling-zstd_test.go)
- [pb-marshaling-only_test.go](bench/pb-marshaling-only_test.go)
- [pb-marshaling-gzip_test.go](bench/pb-marshaling-gzip_test.go)
- [pb-marshaling-lz4_test.go](bench/pb-marshaling-lz4_test.go)
- [pb-marshaling-snappy_test.go](bench/pb-marshaling-snappy_test.go)
- [pb-marshaling-zstd_test.go](bench/pb-marshaling-zstd_test.go)

## Result

같은 파일을 가지고 여러 압축 방식을 통해 Kafka에 저장하고, 이에 따른 저장된 사이즈, 소요시간 등과 같은 리소스 사용량을 확인

> ![topic-list](https://user-images.githubusercontent.com/31076511/142731753-d75b8dc0-1320-477c-9a68-4ab2369e857a.png)

### Size, Duration, CPU/Traffic Usage

| Compression Codec | Size(GB) | Duration(sec)  | Producer CPU(Avg) | Kafka CPU(Avg)  | Kafka Traffic(Avg/Mib)  |
| -------------- |:-----|:---------|:-----|:----|:------|
| Json(Uncomp)   | 22.7 | 432.863  | 25.8 | 9.9 | 411.8 |
| Json + Gzip    | 3.3  | 605.381  | 32.8 | 4.5 | 43    |
| Json + Lz4     | 4.6  | 326.519  | 34.6 | 8.8 | 111.8 |
| Json + Snappy  | 5.5  | 345.558  | 29   | 5.9 | 124.6 |
| Json + Zstd    | 3    | 280.466  | 37.5 | 8.1 | 86.3  |
| ProtoBuffer    | 8.5  | 272.288  | 29.1 | 5.6 | 246.7 |
| PB + Gzip      | 2.8  | 456.264  | 35.4 | 3.2 | 48.8  |
| PB + Lz4       | 3.8  | 244.481  | 34.6 | 6.6 | 122.9 |
| PB + Snappy    | 3.9  | 268.386  | 30.3 | 4.9 | 116.8 |
| PB + Zstd      | 2.6  | 220.791  | 37.2 | 8.2 | 94.8  |

### Diff Ratio
```text
Json 데이터는 uncompress json를 기준으로 비교
ProtoBuffer 데이터는 ProtoBuffer를 기준으로 비교
```

| Compression Codec | Duration Ratio | Producer CPU Ratio | Kafka CPU Ratio | Traffic Ratio | Ratio SUM|
| -------------- |:-----:|:----:|:----:|:----:|:----:|
| Json(Uncomp)   | 1.0 | 1.0 | 1.0 | 1.0 | `4.0` |
| Json + Gzip    | 1.4 | 1.3 | 0.5 | 0.1 | 3.2 |
| Json + Lz4     | 0.8 | 1.3 | 0.9 | 0.3 | 3.3 |
| Json + Snappy  | 0.8 | 1.1 | 0.6 | 0.3 | `2.8` |
| Json + Zstd    | 0.6 | 1.5 | 0.8 | 0.2 | 3.1 |
| ProtoBuffer    | 1.0 | 1.0 | 1.0 | 1.0 | `4.0` |
| PB + Gzip      | 1.7 | 1.2 | 0.6 | 0.2 | 3.7 |
| PB + Lz4       | 0.9 | 1.2 | 1.2 | 0.5 | 3.8 |
| PB + Snappy    | 1.0 | 1.0 | 0.9 | 0.5 | `3.4` |
| PB + Zstd      | 0.8 | 1.3 | 1.5 | 0.4 | 3.9 |

> 다음과 같이 Ratio값을 합산하여 결과를 보는 것이 맞는지는 의문이지만, 어느정도 결과를 도출하기 위해서 사용했다

Ratio 값을 합산하면 Json/ProtoBuffer 두 기준에서 모두 Snappy 압축이 제일 좋은 성능을 가지는 것으로 보인다

실제 적용 시에는 Duration, CPU, Traffic등과 같이 중요시 되는 기준에 따라 Compression Type에 따라 선택한다  



## Options: Resource Usage Graph 
Producer에서 Kafka로 데이터를 저장할때의 compression type별 리소스 사용량이다

### 1. Json Uncompressed
![json-uncompressed](https://user-images.githubusercontent.com/31076511/142730639-f71ac7fc-f7f7-411b-836f-4eef0d337b7e.png)


### 2. Json + Gzip
![json-gzip](https://user-images.githubusercontent.com/31076511/142730649-a56380dd-c616-474c-9c23-f6f6d6ae608c.png)


### 3. Json + Lz4
![json-lz4](https://user-images.githubusercontent.com/31076511/142730664-b68cc57b-733c-4a54-b418-f538a26e98ce.png)


### 4. Json + Snappy
![json-snappy](https://user-images.githubusercontent.com/31076511/142730677-746a5862-3c94-4b44-a6b0-1225278d12d6.png)

   
### 5. Json + Zstd
![json-zstd](https://user-images.githubusercontent.com/31076511/142730682-57714107-6323-4ea8-9324-ee16ad2e86d1.png)


### 6. Protocol Buffer
![protobuf](https://user-images.githubusercontent.com/31076511/142730705-33e1bb4b-e4c9-406e-a583-3487ab9aabc8.png)


### 7. ProtoBuf + Gzip
![pb-gzip](https://user-images.githubusercontent.com/31076511/142730708-64304d71-f593-4a06-ab83-d0d4b2be8d84.png)


### 8. ProtoBuf+ Lz4
![pb-lz4](https://user-images.githubusercontent.com/31076511/142730728-2fffb7e6-9d54-4ca6-84a9-563c174f1075.png)


### 9. ProtoBuf + Snappy
![pb-snappy](https://user-images.githubusercontent.com/31076511/142730731-2c04faa0-5234-4f48-be19-31b959e9fdd3.png)


### 10. ProtoBuf + Zstd
![pb-zstd](https://user-images.githubusercontent.com/31076511/142730740-79101cde-bf95-42c6-b804-ca13256e9f13.png)

## Reference
 - [benefits-compression-kafka-messaging](https://developer.ibm.com/articles/benefits-compression-kafka-messaging/)
 - [KIP-390:A Support Compression Level](https://cwiki.apache.org/confluence/display/KAFKA/KIP-390%3A+Support+Compression+Level)