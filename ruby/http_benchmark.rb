require 'net/http'

raise ArgumentError, 'No URL provided' unless ARGV[0]

URL = URI(ARGV[0]).freeze
results = []

trap "SIGINT" do
  print_results(results)
  exit 130
end

def print_results(results)
  if results.size == 0
    puts 'Nothing to print...'
    return
  end

  puts "Bechmarking HTTP GET to #{URL}"
  puts "avg #{results.sum(0.0) / results.size}"
  puts "min #{results.min}"
  puts "max #{results.max}"
  puts "#{results.count} requests total"
end

while true do
  time = Time.now
  response = Net::HTTP.get_response(URL)
  response.body
  diff = Time.now - time
  puts "Fetched in #{diff} seconds"
  results << diff
end
