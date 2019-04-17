#^ sip:219.142.69.234:9080
#^ sip:172.17.0.2:9200
#^ sip:s.berrybit.cn
#^ sip:wxvision.ruibei365.com
#^ sip:127.0.0.1:8090

GET ' ' http://127.0.0.1:8880/msb/services

GET ' ' https://127.0.0.1:8880/

GET ' ' https://127.0.0.1:8880/msa/healthcheck

GET ' ' https://127.0.0.1:8880/


POST '
{
    "ServiceName": "ServiceName",
	"Version": "Version",
	"Desc": "Desc",
	"IPAddr": "IPAddr",
	"Port": 9999999,
	"HostName": "HostName",
	"ProjName": "ProjName",
	"CreatedAt": "CreatedAt"
}
' http://ip/service

GET ' ' http://ip/service/ALL/Version
DELETE ' ' http://ip/service/ALL/Version


POST '
{

	"DevType": "XXX",
	"DevGUID": "XXX",
	"DevVersion": "XXX",
	"UserGUID": "XXX",
	"TimeStamp": "XXX",
	"LeftAxis1": "XXX",
	"LeftAxis2": "XXX",
	"LeftDc1": "XXX",
	"LeftDc2": "XXX",
	"LeftDs1": "XXX",
	"LeftDs2": "XXX",
	"LeftGazeh": "XXX",
	"LeftGazev": "XXX",
	"LeftPupil": "XXX",
	"LeftSe1": "XXX",
	"LeftSe2": "XXX",
	"RightAxis1": "XXX",
	"RightAxis2": "XXX",
	"RightDc1": "XXX",
	"RightDc2": "XXX",
	"RightDs1": "XXX",
	"RightDs2": "XXX",
	"RightGazeh": "XXX",
	"RightGazev": "XXX",
	"RightPupil": "XXX",
	"RightSe1": "XXX",
	"RightSe2": "XXX"
}
' http://ip/check/diop

POST '
{
	"DevType": "XXX",
	"DevGUID": "XXX",
	"DevVersion": "XXX",
	"UserGUID": "XXX",
	"TimeStamp": "XXX",
	"Left": "XXX",
	"Right": "XXX"
}
' http://ip/check/vision

