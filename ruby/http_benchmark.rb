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
  durations = results.map(&:first)

  puts "\nBechmarking HTTP GET to #{URL}"
  puts "avg #{durations.sum(0.0) / durations.size}"
  puts "min #{durations.min}"
  puts "max #{durations.max}"
  puts "#{results.count} requests total"
end

while true do
  time = Time.now
  response = Net::HTTP.get_response(URL)
  response.body
  diff = Time.now - time
  puts "fetched in #{diff} seconds, code: #{response.code}"
  results << [diff, response.code]
end
