# frozen_string_literal: true

require 'socket'
require 'securerandom'

port = 8001

s = TCPSocket.open('127.0.0.1', port)
buffer = ''

puts ':::Start receive binary data'
binary = s.gets
until binary == 'EOF'
  binary = s.gets
  break if binary.nil?

  buffer += binary
end

s.close
puts ':::Received binary data'
file_name = "#{SecureRandom.uuid}_received_se.mp3"
File.open(file_name, 'w') do |f|
  f.write(buffer)
end

puts ':::Play sound'
system('afplay', file_name)
File.delete(file_name)

puts 'finished!!!'
