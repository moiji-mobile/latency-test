#!/usr/bin/ruby
require 'socket'

port = ARGV.first.to_i
puts "S: Listening on #{port}"
server = TCPServer.open(port)


while true
    client = server.accept
    pkts = 0
    while true
        hdr = client.read(2)
        if not hdr
            puts "S: short read #{pkts}"
            break
        end
        if hdr.length != 2
            puts "S: short header #{hdr.length} #{pkts}"
            break
        end
        len = hdr.unpack('n')[0]
        dat = client.read(len)
        if dat.length != len
            puts "S: short read #{hdr.length} #{len} #{pkts}"
            break
        end
        rc = client.write hdr+dat
        pkts += 1
    end
    client.close
end
