require 'sinatra'
require 'net/http'

set :bind, '0.0.0.0'
set :port, 3000

get '/recommend' do
  content_type :json
  user = request.env["HTTP_USER"]
  puts "get recommend products for user #{user}"

  ids = getRecommendIDsByUser(user)
  products = getProducts(ids)
  products.to_json
end

def getRecommendIDsByUser(user)
  mockDB = (1..4).to_a
  #mockDB = (7..20).to_a
  mockDB.sample(4)
end

def getProducts(ids)
  url = URI("http://127.0.0.1:5000/products?ids=#{ids.join(',')}")
  response = Net::HTTP.get(url)
  JSON.parse(response)
end

