package bench

import (
	"encoding/csv"
	"github.com/99-66/compression-efficiency-in-kafka/models"
	pb "github.com/99-66/compression-efficiency-in-kafka/protos/music"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
)

// mapToMusic CSV record를 모델로 매핑하여 모델로 반환한다
func mapToMusic(record []string) models.Music {
	music := models.Music{
		BuySeq:         record[0],
		ProductCd:      record[1],
		BuyId:          record[2],
		ProductSeq:     record[3],
		ProductNM:      record[4],
		BuyNM:          record[5],
		BeforeMyMoney:  record[6],
		SendStatus:     record[7],
		ProductCnt:     record[8],
		SendCnt:        record[9],
		DenyCnt:        record[10],
		ReceiveCnt:     record[11],
		SendTm:         record[12],
		ProductMyMoney: record[13],
		BuyMyMoney:     record[14],
		DiscountRate:   record[15],
		BuyIpAddr:      record[16],
		BuyLoginId:     record[17],
		BuyLoginName:   record[18],
		ReceiveTm:      record[19],
	}

	return music
}

// mapToMusic CSV record를 프로토콜버퍼 모델로 매핑하여 모델로 반환한다
func mapToProtoBufferMusic(record []string) *pb.Music {
	music := pb.Music{
		BuySeq:         record[0],
		ProductCd:      record[1],
		BuyId:          record[2],
		ProductSeq:     record[3],
		ProductNm:      record[4],
		BuyNm:          record[5],
		BeforeMyMoney:  record[6],
		SendStatus:     record[7],
		ProductCnt:     record[8],
		SendCnt:        record[9],
		DenyCnt:        record[10],
		ReceiveCnt:     record[11],
		SendTm:         record[12],
		ProductMyMoney: record[13],
		BuyMyMoney:     record[14],
		DiscountRate:   record[15],
		BuyIpAddr:      record[16],
		BuyLoginId:     record[17],
		BuyLoginName:   record[18],
		ReceiveTm:      record[19],
	}

	return &music
}

func generateJsonSample(r io.Reader) (ch chan models.Music) {
	ch = make(chan models.Music)
	go func() {
		reader := csv.NewReader(r)
		if _, err := reader.Read(); err != nil {
			log.Fatalln(err)
		}
		defer close(ch)
		for {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalln(err)
			}
			music := mapToMusic(record)
			ch <- music
		}
	}()

	return
}

// MusicData Sample 데이터를 전달하기 위한 채널용 모델을 선언한다
type MusicData struct {
	Bytes []byte
	Err   error
}

func generateProtoBufferSample(r io.Reader) (ch chan MusicData) {
	ch = make(chan MusicData)
	go func() {
		reader := csv.NewReader(r)
		if _, err := reader.Read(); err != nil {
			log.Fatalln(err)
		}
		defer close(ch)
		for {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalln(err)
			}
			music := mapToProtoBufferMusic(record)
			encodedMusic, err := proto.Marshal(music)
			if err != nil {
				ch <- MusicData{
					Bytes: nil,
					Err:   err,
				}
				log.Printf("generated failed...%s\n", record)
			}

			ch <- MusicData{
				Bytes: encodedMusic,
				Err:   nil,
			}
		}
	}()

	return
}

var samples = []string{
	"data_10k",
	"data_1m",
	"data_5m",
	"data_10m",
	"data_50m",
	"data_100m",
}
