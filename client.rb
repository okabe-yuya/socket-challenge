require 'socket'

1.upto(5) do
  socket = TCPSocket.open('0.0.0.0', 8001)
  puts socket.gets
  sleep 1
end
