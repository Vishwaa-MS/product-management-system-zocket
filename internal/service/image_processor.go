package service

import (
	"log"
	"product-management-system/internal/queue"
	"time"
)

// ImageProcessor handles asynchronous image processing tasks
type ImageProcessor struct {
	Queue *queue.RabbitMQ
}

// NewImageProcessor creates a new ImageProcessor instance
func NewImageProcessor(queue *queue.RabbitMQ) *ImageProcessor {
	return &ImageProcessor{
		Queue: queue,
	}
}

// ConsumeImageProcessingQueue starts consuming messages from the queue
func (p *ImageProcessor) ConsumeImageProcessingQueue() {
	go func() {
		for {
			msg, ok := <-p.Queue.ConsumeMessages()
			if !ok {
				log.Println("Queue closed or stopped")
				break
			}

			// Simulate image processing (placeholder logic)
			log.Printf("Processing image: %s", msg)
			time.Sleep(2 * time.Second) // Simulate processing time
		}
	}()
}
