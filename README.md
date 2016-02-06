# gotification

## usage

1. gotification.Option{GCM string, APN string}.Set() **GCM is apikey, APN is filename of certification file.**
2. gotification.Notification{Message string, IOSReceivers []string, AndroidReceivers []string}
3. gotification.Notification.Send()
