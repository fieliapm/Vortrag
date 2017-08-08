(info) after VOEZ launch - how to resolve problems of mobile game server development and service maintenance
===

# Info

## Abstract

As startup game makers, many people might be able to build a workable mobile game service, but have no experience in building a stable, reliable, high performance commercial mobile game service.
Sometimes we feel that our code seems to be logical correct and should work, but in reality it failed because we didn't avoid race condition very well, or we didn't apply caching mechanism appropriately.
In this talk, we will introduce how we built a game server, which is fully implemented in Python and Flask.
We are going to share some experiences with backend or online mobile game developers by showing the essence of Python code and Flask usage from game server.
Hope these guides will lead you to build a stable, reliable service in Python, and don't repeat the same mistakes we did before.

## Domain Knowledge:

	RSA or HMAC principle and usage
	HTTP protocol:
		error code
		cache mechanism
	AWS or GCP experience

## Target

To share experience and Python code to newbie game developers who want to build stable, reliable, high performance mobile game service.
To encourage senior game developers share better implementation strategy than us.

## Detailed description

What we will share at below chapter:

VOEZ game play demo and connection layout:
1

Python code:
2-1 game play authorization and score uploading
3-2 database/storage layout & API for distributing current seasonal event revision and corresponding assets
3-3 Redis database operation
3-4 order of inter-server request and related state transition
4-2 publish event game data to Amazon S3/Google Cloud Storage
4-3 append header "Cache-Control" under Flask framework, and make cache expire precisely at seasonal event switching
5-1 calculate how many consecutive days a player logged in

Cloud platform setting tips:
4-1 What happened when we met DDoS

========================================

# Chapter

1. VOEZ current status [3min]
2. HTTP protocol [1min]
3. genuine & purchase verification [5min]
    1. clean leaderboard: signature
    2. legal game play: activation
4. stability [7min]
    1. principle
    2. publish order and request order
    3. database cache mechanism
    4. reliable inter-server request
    5. execution resource allocation for request handler and for database
5. performance [7min]
    1. from database to static file
    2. statistics
    3. server operation with CDN
    4. service downtime and update
6. timezone [2min]
    1. all about environment variable: TZ
7. conclusion [1min]
8. Q&A [3min]

========================================

# Detail

## genuine and purchase verification

### request & response signature

	score:
		request: use RSA encrypt message or message checksum plus timestamp and nonce (a random variable)
	activation:
		request: use RSA encrypt message or message checksum plus timestamp and nonce (a random variable)
		response: use any signature algorithm, nonce from encrypted message as key, to sign response data

## stability

### principle

	always read and write game data to DB atomically (DB status should be before request and after request, there is no DB status in the middle of request)

### publish order and request order

	if: client read game info A, then game info B, then game info C
	then: publish game info C, then B, then A
	database write and read schema can also use the above principle

### database cache mechanism

	read cached data
	1. try to read cache
	2. if cache expired or not exist, read original data
	3. update to cache and return data
	* caution: a data should be set or get atomically
	* caution: don't build cache content based on previous cache content
        *  ex: increase/drecrease counter

### reliable inter-server request

	inter server connection and transaction:
	server A calculate data
	server A start transaction and write data
	server A request server B
	if success:
		server A close transaction
	else:
		server A revert transaction

## performance

### execution resource allocation for request handler and for database

	database processing speed per request should be superior to request handler (so we can queue huge requests with load balancer but keep database server reliable)

### from database to static file

	make all updatable game info become static file and feed to client
	if you can implement with static file, then do not implement with database

### some implementation detail while operating with CDN

	asset update
		different revision of file (ex. content changed) should locate in different http resource path
			add checksum in http resource path
		event switch control
			server expire time via Cache-Control
			client can use this expire time as refresh timer (back to title screen and update event)

	export log to big query data
		alert maintainer
		data analysis
		seasonal event result
        
### service downtime and update

    cachable: GET 200, 203, 300, 301, 302, 307, 410 (when using Google CDN)
    downtime: 503 (non-cachable: even if CDN exists, success HTTP response will be sent to client as long as server become alive again)
    update: 410 (cachable: CDN will tell client "this path is abandoned" without bothering original server)

## timezone

### all about environment variable: TZ

if you want to announce event and calculate login day accumulation belong to localtime, set TZ to timezone to what you want to refer.
If daylight saving time is applied in this timezone, then UTC offset of this timezone will vary, and this variation will automatically follow rules of this timezone or country.

