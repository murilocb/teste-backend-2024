class ApplicationConsumer < Karafka::BaseConsumer
  def consume
    messages.each do |message|
      data = message.payload
      Rails.logger.info("Processing message with data: #{data}")

      begin
        ActiveRecord::Base.transaction do
          product = Product.find_or_initialize_by(id: data["id"])
          product.assign_attributes(
            name: data["name"],
            description: data["description"],
            brand: data["brand"],
            price: data["price"],
            stock: data["stock"],
            created_at: parse_timestamp(data["created_at"]),
            updated_at: parse_timestamp(data["updated_at"])
          )
          product.save!
          Rails.logger.info("Product #{product.new_record? ? 'created' : 'updated'} with ID: #{product.id}")
        end
      rescue => e
        Rails.logger.error("Failed to process message: #{e.message}")
        raise e
      end
    end
  end

  private

  def parse_timestamp(timestamp_str)
    Time.iso8601(timestamp_str) rescue nil
  end
end
