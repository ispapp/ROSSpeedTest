# ROSSpeedTest

### Transmission (TX) Formula

When transmitting a data packet over a communication channel, the time it takes (\(TX\) time) can be calculated using the formula:

> $TX \, \text{Time} = \frac{\text{Packet Size}}{\text{Transmission Speed}}$

Where:
- **Packet Size:** The size of the data packet being transmitted, measured in bits.
- **Transmission Speed:** The rate at which data can be transmitted on the communication channel, measured in bits per second.

### Reception (RX) Formula

When receiving a data packet over a communication channel, the time it takes (\(RX\) time) can be calculated using the formula:

> $RX \, \text{Time} = \frac{\text{Packet Size}}{\text{Reception Speed}}$

Where:
- **Packet Size:** The size of the data packet being received, measured in bits.
- **Reception Speed:** The rate at which data can be received on the communication channel, measured in bits per second.

### Ping (Latency) Formula

The total round-trip time (RTT) or latency for a small piece of data to travel from the source to the destination and back (\(Ping\)) can be calculated using the formula:

> $Ping \, (\text{Latency}) = \text{TX Time} + \text{Propagation Delay} + \text{Processing Delay} + \text{Queuing Delay} + \text{RX Time}$

Where:
- **TX Time:** Transmission time for sending data from the source to the destination.
- **Propagation Delay:** Time for the signal to travel through the communication medium.
- **Processing Delay:** Delay within devices (routers, switches) along the communication path.
- **Queuing Delay:** Delay due to waiting in a queue in a congested network.
- **RX Time:** Reception time for receiving data from the destination back to the source.
