class ApplicationJob < ActiveJob::Base
  queue_as :default

  def perform(product_params)
    kafka_producer = KafkaProducerService.new
    kafka_producer.publish("rails-to-go", product_params.to_json)
  end
end
