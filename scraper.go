package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/giang19062001/chi-golang/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// TODO: 1 VÒNG LẶP VÔ HẠN CỨ MỖI 1 PHÚT -> LẤY 10 FEED TỪ DATABASE -> CÁC GOROUTINE GỌI URL LẤY XML DATA VÀ CHẠY LOG "TITLE" TỪNG "ITEM" CỦA TỪNG FEED_ID
// TODO: -> WaitGroup chờ tất cả xong -> Quay lại vòng lặp chờ đến 1 phút tiếp theo

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetweenRequest, concurrency)

	// tạo ra một ticker — là một “đồng hồ” sẽ gửi tín hiệu đều đặn mỗi 'timeBetweenRequest' giây
	ticker := time.NewTicker(timeBetweenRequest)

	// vòng lặp vô hạn
	// sau mỗi lần lặp -> chương trình dừng lại chờ đến khi 'ticker.C' gửi tín hiệu tiếp theo
	for ; ; <-ticker.C {
		// ticker.C là một "channel" kiểu <-chan time.Time
		// ! ko giống như r.Context() từ request http, đây là chạy ngầm nên context phải là context.Background()
		feeds, err := db.GetNextFeedToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}
		log.Printf("***** Found %v feeds to fetch *****", len(feeds))

		// Tạo WaitGroup để chờ các goroutine hoàn thành
		wg := &sync.WaitGroup{}
		// lặp qua danh sách các feed từ database ( hiện tại tối đa 10 )
		for _, feed := range feeds {
			// báo rằng có thêm 1 goroutine cần chờ
			wg.Add(1)
			// ! gọi go routine xử lý song song  => có nghĩa 10 feed có thể chạy hàm này cùng lúc
			go scrapeFeed(db, wg, feed)
		}
		// chờ cho tất cả goroutines hoàn tất trước khi quay lại vòng lặp
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	// Dù thành công hay lỗi, khi hàm kết thúc, WaitGroup giảm đếm xuống 1
	defer wg.Done()
	// cập nhập cột last_fetched_At bằng thời gian hiện tại
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	// gọi URL lấy danh sách xml data của feed hiện tại
	feedData, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	// log từng tiêu đề của item của feed data theo cấu trúc struct đã khai báo
	for _, item := range feedData.Channel.Item {
		log.Println("Found post ==>: ", item.Title)
	}
	log.Printf("------- Feed %s collected, %v posts found -------", feed.Name, len(feedData.Channel.Item))
}

func fetchFeed(feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		// nếu server không phản hồi trong 10 giây, request sẽ bị hủy → tránh treo chương trình
		Timeout: 10 * time.Second,
	}

	// GET để tải nội dung từ URL
	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // đảm bảo đóng kết nối khi hàm kết thúc (dù lỗi hay không)

	log.Println("resp.Body", resp.Body)
	// Đọc toàn bộ dữ liệu XML từ response
	dat, err := io.ReadAll(resp.Body)
	// log.Println("dat", dat)
	if err != nil {
		return nil, err
	}

	// chuyển dữ liệu XML → struct.
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}
	return &rssFeed, nil

}
