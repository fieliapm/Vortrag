after VOEZ launch
===

<!-- .slide: data-background="#FFDFEF" -->
<!-- .slide: data-transition="zoom" -->

how to resolve problems of mobile game server development and service maintenance :dizzy:

> [name=郭學聰 Hsueh-Tsung Kuo] [time=Fri, 09 Jun 2017] [color=red]

---

<!-- .slide: data-transition="convex" -->

## who am I?

![fieliapm](https://pbs.twimg.com/profile_images/591670980021387264/aZAYLRUe_400x400.png)

----

<!-- .slide: data-transition="convex" -->

* programmer from Rayark, a game company in Taiwan
* backend engineer
* usually develop something related to my work in Python, Ruby, Golang, C#
* built almost entire VOEZ game server by myself only

---

<!-- .slide: data-transition="convex" -->

## VOEZ current status

----

### spec

<!-- .slide: data-transition="convex" -->

* rhythm game following after famous titles Cytus and Deemo
* profile, save data, achievement, leaderboard... are stored on server
* based on Google Cloud Platform

----

<!-- .slide: data-transition="convex" -->

### resource

* official website: https://www.rayark.com/g/voez/
* teaser: https://youtu.be/Bh6gQyJHbxI
* walkthrough wiki: http://voez.info/

* <small>Apple iTunes: https://itunes.apple.com/jp/app/voez/id1007929736</small>
* <small>Google Play: https://play.google.com/store/apps/details?id=com.rayark.valkyrie</small>
* <small>Nintendo Switch: https://ec.nintendo.com/JP/ja/titles/70010000000044</small>

----

<!-- .slide: data-transition="convex" -->

### teaser

{%youtube Bh6gQyJHbxI %}

---

<!-- .slide: data-transition="convex" -->

## outline

----

<!-- .slide: data-transition="convex" -->

3. VOEZ current status

----

<!-- .slide: data-transition="convex" -->

5. HTTP protocol
6. genuine & purchase verification
    1. clean leaderboard: signature
    2. legal game play: activation
7. stability
	1. principle
	2. publish order and request order
	3. database cache mechanism
	4. reliable inter-server request
	5. execution resource allocation for request handler and for database

----

<!-- .slide: data-transition="convex" -->

8. performance
	1. from database to static file
	2. statistics
	3. server operation with CDN
    4. service downtime and update
9. timezone
    1. all about environment variable: TZ
10. conclusion
11. Q&A

---

<!-- .slide: data-transition="convex" -->

## HTTP protocol

----

<!-- .slide: data-transition="convex" -->

request

```
POST /api/leaderboard/song/add_score HTTP/1.1
Host: voez-api.rayark.net
Signature: cHl0aG9uaXN0YQ==

{
  "token": "593a19bc593a19bc",
  "player": "pythonista",
  "timestamp": 1496979900,
  "score": {
    "total": 999999.0,
    ......
  }
}
```

----

<!-- .slide: data-transition="convex" -->

response

```
HTTP/1.1 200 OK
Server: nginx/1.9.12
Content-Type: application/json; charset=UTF-8
Content-Length: 22
Cache-Control: public,s-maxage=300
Date: Fri, 09 Jun 2017 10:20:30 GMT

{
    "status": "ok"
}
```

----

<!-- .slide: data-transition="convex" -->

API layout

1. gamedata/event/asset download
2. player info setting/getting
3. leaderboard (top/self/friends)
4. start-play authorization & score uploading
5. avatar lottery

---

<!-- .slide: data-transition="convex" -->

## genuine
## &
## purchase verification

----

<!-- .slide: data-transition="convex" -->

* ~~leaderboard~~ <- fake score
* ~~game play~~ <- pirated app & private server

----

<!-- .slide: data-transition="convex" -->

### clean leaderboard: signature

- request: send request & attach ==*signature*== <!-- .element: class="fragment" data-fragment-index="1" -->

----

<!-- .slide: data-transition="convex" -->

```sequence
client->server: request w/ signature
note right of server: verify signature
server->client: status
```

----

<!-- .slide: data-transition="convex" -->

```python=
# client side
def sign_message(private_key_string, message):
    private_key = Crypto.PublicKey.RSA.importKey(private_key_string)
    signer = Crypto.Signature.PKCS1_v1_5.new(private_key)
    message_hash = Crypto.Hash.SHA256.new(message)
    signature = signer.sign(message_hash)
    return base64.standard_b64encode(signature).decode('utf-8')

# server side
def verify_message(public_key_string, message, signature_b64):
    public_key = Crypto.PublicKey.RSA.importKey(public_key_string)
    verifier = Crypto.Signature.PKCS1_v1_5.new(public_key)
    message_hash = Crypto.Hash.SHA256.new(message)
    signature = base64.standard_b64decode(signature_b64.encode('utf-8'))
    return verifier.verify(message_hash, signature)
```

----

<!-- .slide: data-transition="convex" -->

### legal game play: activation

* request: send request & RSA encrypted request hash plus **nonce** (a random variable)
* response: send response & attach signature
  * ==sign response hash plus **nonce**== <!-- .element: class="fragment" data-fragment-index="1" -->
  * ==client verify signature== <!-- .element: class="fragment" data-fragment-index="2" -->

----

<!-- .slide: data-transition="convex" -->

```sequence
client->server: request w/ RSA.encrypt(merge(request.hash, nonce))
note right of server: verify request.hash
server->client: response w/ RSA.sign(merge(response.hash, nonce))
note left of client: verify signature
```

----

<!-- .slide: data-transition="convex" -->

```python=
# client side
def encrypt_message(public_key_string, message):
    public_key = Crypto.PublicKey.RSA.importKey(public_key_string)
    cipher = Crypto.Cipher.PKCS1_OAEP.new(public_key)
    encrypted_message = cipher.encrypt(message)
    return base64.standard_b64encode(encrypted_message).decode('utf-8')

def encrypt_message_hash_with_nonce(public_key_string, message, nonce):
    message_hash_digest = Crypto.Hash.SHA256.new(message).digest() # this computation may take too long
    return encrypt_message(public_key_string, merge(message_hash_digest, nonce))

# server side
def decrypt_message(private_key_string, encrypted_message_b64):
    private_key = Crypto.PublicKey.RSA.importKey(private_key_string)
    cipher = Crypto.Cipher.PKCS1_OAEP.new(private_key)
    encrypted_message = base64.standard_b64decode(encrypted_message_b64.encode('utf-8'))
    return cipher.decrypt(encrypted_message)

def verify_encrypted_message_hash_and_extract_nonce(private_key_string, message, encrypted_message_hash_with_nonce_b64):
    message_hash_with_nonce = decrypt_message(private_key_string, encrypted_message_hash_with_nonce_b64)
    message_hash_digest = Crypto.Hash.SHA256.new(message).digest() # this computation may take too long
    (decrypted_message_hash_digest, nonce) = unmerge(message_hash_with_nonce)
    return (decrypted_message_hash_digest == message_hash_digest, nonce)
```

---

<!-- .slide: data-transition="convex" -->

## stability

----

<!-- .slide: data-transition="convex" -->

### principle

* always read & write game data to DB atomically
  * available DB updating status
    * before request:
      * start transaction, not update yet
    * after request:
      * success: update completed
      * fail: not update
* caution:
  * impossible to follow this rule everywhere

----

<!-- .slide: data-transition="convex" -->

### publish order and request order

* request order
  * client read game info A, then B, then C
* publish order <!-- .element: class="fragment" data-fragment-index="1" -->
  * publish game info C, then B, then A
  * remove game info A, then B, then C
* the above principle is suitable for database and online resource distribution <!-- .element: class="fragment" data-fragment-index="2" -->

----

<!-- .slide: data-transition="convex" -->

```python=
class Revision(mongo_engine.Document): # read this 1st, publish this 2nd, remove this 1st
    revision_id = mongo_engine.StringField(required=True)
    timestamp = mongo_engine.DateTimeField(required=True)
    name = mongo_engine.StringField(required=True)
    meta = {
        'indexes': [
            {
                'fields': ['timestamp'],
                'unique': True,
            },
        ],
    }

class SongAssetMeta(mongo_engine.Document): # read this 2nd, publish this 1st, remove this 2nd
    revision_id = mongo_engine.StringField(required=True)
    song_id = mongo_engine.StringField(required=True)
    asset_set_checksum = mongo_engine.StringField(required=True)
    song_cls = mongo_engine.StringField(required=True)
    song_pack_id = mongo_engine.StringField(required=False)
    song_title = mongo_engine.StringField(required=True)
    meta = {
        'indexes': [
            {
                'fields': ['revision_id', 'song_id'],
                'unique': True,
            },
        ]
    }
```

----

<!-- .slide: data-transition="convex" -->

### database cache mechanism

* database in HDD or SSD: slow
* Redis or Memcached in RAM: fast

----

<!-- .slide: data-transition="convex" -->

* workable mechanism in concurrency environment
  1. try to read cache
  2. if cache missed, read original data and update data to cache 
  3. read cache and return data
* caution
  * cache should be set or get atomically
  * don't build cache from previous cache content
    * ex: increase & drecrease counter

----

<!-- .slide: data-transition="convex" -->

```python=
class SongAssetMetaCache(object):
    def set_meta(self, revision_id, song_id, asset_set_checksum, song_cls, song_pack_id):
        meta_key = 'song_asset_meta:%s:%s' % (revision_id, song_id)
        song_meta = {'asset_set_checksum': asset_set_checksum, 'song_cls': song_cls, 'song_pack_id': song_pack_id}
        self.strict_redis.hmset(meta_key, song_meta)
        self.strict_redis.expire(meta_key, self.cache_expire_time)

    def get_meta(self, revision_id, song_id):
        return self.strict_redis.hmget(revision_id, song_id)

class SongAssetMetaModel(object):
    def reconstruct_cache(self, revision_id, song_id):
        try:
            song_asset_meta = SongAssetMeta.objects.get(revision_id=revision_id, song_id=song_id)
        except SongAssetMeta.DoesNotExist:
            SongAssetMeta.delete(revision_id, song_id)
        else:
            self.__song_asset_meta_cache.set_meta(song_asset_meta.revision_id, song_asset_meta.song_id,
                song_asset_meta.asset_set_checksum, song_asset_meta.song_cls, song_asset_meta.song_pack_id)

    def get_song_asset_meta(self, revision_id, song_id):
        song_asset_meta = self.__song_asset_meta_cache.get_meta(revision_id, song_id)
        if song_asset_meta is None:
            self.reconstruct_cache(revision_id, song_id)
            song_asset_meta = self.__song_asset_meta_cache.get_meta(revision_id, song_id)
        return song_asset_meta
```

----

<!-- .slide: data-transition="convex" -->

### reliable inter-server request

* inter server connection and transaction:
  1. server A calculate data
  2. server A start transaction and write data
  3. server A request server B
     * success: server A finish transaction
     * fail: server A revert transaction

----

<!-- .slide: data-transition="convex" -->

```python=
def consume(player_access_token):
    failure_count = 0
    while True:
        try:
            purchase.consume_coin(player_access_token)
        except (requests.exceptions.RequestException, server_error.ServerError):
            failure_count += 1
            if failure_count >= 3:
                raise
        else:
            break

def gacha(player_access_token):
    selected_avatars = randomly_select_avatars(player_access_token)
    transaction = begin_appending_avatars_to_player_data(player_access_token, selected_avatars)
    try:
        consume(player_access_token)
    except:
        revert_appending_avatars_to_player_data(transaction)
        raise
    else:
        finish_appending_avatars_to_player_data(transaction)
```

----

<!-- .slide: data-transition="convex" -->

### execution resource allocation for request handler and for database

* database processing speed per request should be superior to request handler
  * [x] queue huge requests with load balancer but keep database server reliable

----

<!-- .slide: data-transition="convex" -->

![SINoALICE](https://i.imgur.com/0LZjGBO.jpg)
https://i.imgur.com/0LZjGBO.jpg

---

<!-- .slide: data-transition="convex" -->

## performance

----

<!-- .slide: data-transition="convex" -->

#### Python is so SLOW

> :hash: "skip Python code execution
> if you can" 
> [name=Hsueh-Tsung Kuo] [time=Fri, 09 Jun 2017] [color=red]
> <!-- .element: class="fragment" data-fragment-index="1" -->

----

<!-- .slide: data-transition="convex" -->

### from database to static file

* announce public game info from static file
* only per-user data is served from database

----

<!-- .slide: data-transition="convex" -->

```python=
def publish_asset_revision(revision):
    song_list = [process(song) for song in SongAssetMeta.objects(revision_id=revision.id)]
    song_list_json = json.dumps(song_list, ensure_ascii=False, default=bson.json_util.default)
    location = get_cloud_storage_location()
    with io.BytesIO(json_data) as fp:
        cloud_storage.upload_file(location, fp, 'application/json', len(json_data))
```

----

<!-- .slide: data-transition="convex" -->

### statistics

* what can we do when exporting log to Google BigQuery
  * data analysis (boring)
  * find unusual behavior and alert maintainer (help DevOps)
  * collect & announce seasonal event result (!?)

----

<!-- .slide: data-transition="convex" -->

#### collect & announce seasonal event result

* run script with specific intervals (using crontab)
  * run BigQuery and collect result
  * save result to Google Cloud Storage and make it public

----

<!-- .slide: data-transition="convex" -->

### server operation with CDN

```sequence
client->CDN: request 1
Note left of CDN: cache miss
CDN-->server: request 1
server-->CDN: response 1
CDN->client: response 1

client->CDN: request 2
Note right of CDN: cache hit
CDN->client: response 2 (same as response 1)
```

----

<!-- .slide: data-transition="convex" -->

#### headers to operating with CDN

```
Cache-Control: public,s-maxage=300
Date: Fri, 09 Jun 2017 10:20:30 GMT
```

----

<!-- .slide: data-transition="convex" -->

#### tips when operating with CDN

* asset files
  * different revisions of asset files should be located at individual URLs
    * URL must contains revision ID or checksum
* entry data which lists asset files
  * attach event end time via **Cache-Control** or **Expires**
  * CDN will cache contents until event ended
  * client can use expire time as refresh timer
    * ex: back to main menu & display updated game event

----

<!-- .slide: data-transition="convex" -->

```python=
@app.route('/api/asset/asset_info/<directory>/<file_name>', methods=('GET',))
@cache_control(get_cdn_cache_maxage())
def get_asset_info(directory, file_name):
    # handle asset info
    return flask.Response(response=asset_info_json_data, mimetype='application/json')

@app.route('/api/asset/song_list', methods=('GET',)
template_response_headers()
def get_song_list():
    remaining_second = get_current_song_list_remaining_second()
    cache_maxage = min(remaining_second, get_cdn_cache_maxage())
    # handle song list
    response = flask.Response(response=song_list_json_data, mimetype='application/json')
    add_precise_cache_control_to_headers(response.headers, cache_maxage)
    return response
```

----

<!-- .slide: data-transition="convex" -->

```python=
def __add_cache_control_to_headers(headers, s_maxage):
    headers['Cache-Control'] = 'public,s-maxage=%d' % (s_maxage,)

def add_precise_cache_control_to_headers(headers, s_maxage):
    __add_cache_control_to_headers(headers, s_maxage)

def template_response_headers(headers={}):
    def decorator(func):
        @wraps(func)
        def decorated_function(*args, **kwargs):
            flask.g.current_unix_timestamp = time.time()
            response = flask.make_response(func(*args, **kwargs))
            original_headers = response.headers
            for (header, value) in headers.items():
                original_headers.setdefault(header, value)
            original_headers.setdefault('Date', werkzeug.http.http_date(flask.g.current_unix_timestamp))
            return response
        return decorated_function
    return decorator

def cache_control(s_maxage):
    headers = {}
    __add_cache_control_to_headers(headers, s_maxage)
    return template_response_headers(headers)
```

----

<!-- .slide: data-transition="convex" -->

### service downtime and update

* cachable: GET 200, 203, 300, 301, 302, 307, 410 (when using Google CDN)
* downtime: 503 (non-cachable: even if CDN exists, success HTTP response will be sent to client as long as server become alive again)
* update: 410 (cachable: CDN will tell client "this path is abandoned" without bothering original server)

---

<!-- .slide: data-transition="convex" -->

## timezone

----

<!-- .slide: data-transition="convex" -->

### all about environment variable: TZ

:100: the incredible variable:
# TZ=Asia/Taipei <!-- .element: class="fragment" data-fragment-index="1" -->

```python=
os.environ['TZ'] = 'Europe/Berlin'
time.tzset()
```
<!-- .element: class="fragment" data-fragment-index="2" -->

----

<!-- .slide: data-transition="convex" -->

### all about environment variable: TZ

* if you want to announce event and calculate login day accumulation belong to localtime, set TZ to timezone to which you want to refer.
* if daylight saving time is applied in this timezone, then UTC offset of this timezone will vary, and this variation will automatically follow rules of this timezone or country.

----

<!-- .slide: data-transition="convex" -->

```python=
def get_utc_offset_in_second(timestamp):
    struct_time = time.localtime(timestamp)
    return calendar.timegm(struct_time)-int(time.mktime(struct_time))

def get_utc_offset_in_timedelta(timestamp):
    return datetime.datetime.fromtimestamp(timestamp)-datetime.datetime.utcfromtimestamp(timestamp)

def get_unix_epoch_day_in_day(timestamp, utc_offset_in_second, local_day_boundary_in_second=None):
    if local_day_boundary_in_second is None:
        local_day_boundary_in_second = 0
    return (timestamp+utc_offset_in_second-local_day_boundary_in_second)//86400

def get_unix_epoch_day_in_datetime(datetime_obj, utc_offset_in_timedelta, local_day_boundary_in_timedelta=None):
    if local_day_boundary_in_timedelta is None:
        local_day_boundary_in_timedelta = datetime.timedelta()
    return (datetime_obj+utc_offset_in_timedelta-local_day_boundary_in_timedelta).replace(hour=0, minute=0, second=0, microsecond=0)
```

---

<!-- .slide: data-transition="convex" -->

## conclusion

----

<!-- .slide: data-transition="convex" -->

> :hash: "don't repeat the same mistakes we did before!"
> [name=Hsueh-Tsung Kuo] [time=Fri, 09 Jun 2017] [color=red]

----

<!-- .slide: data-transition="convex" -->

### special thanks

* Rayark Inc.
  * CTO & CIO
  * VOEZ team
  * backend team
  * QA team
  * customer service team
  * IT team
  * other teams
* iKala Interactive Media Inc.

---

<!-- .slide: data-transition="zoom" -->

## Q&A

---

<style>

.reveal code {
    font-size: 12px !important;
    line-height: 1.2;
}

body {
    background-color: Indigo;
}

.rightpart{
    float:right;
    width:50%;
}

.leftpart{
    margin-right: 50% !important;
    height:50%;
}
.reveal section img { background:none; border:none; box-shadow:none; }
p.blo {
	font-size: 50px !important;
	background:#B6BDBB;
	border:1px solid silver;
	display:inline-block;
	padding:0.5em 0.75em;
	border-radius: 10px;
	box-shadow: 5px 5px 5px #666;
}

p.blo1 {
	background: #c7c2bb;
}
p.blo2 {
	background: #b8c0c8;
}
p.blo3 {
	background: #c7cedd;
}

p.bloT {
	font-size: 60px !important;
	background:#B6BDD3;
	border:1px solid silver;
	display:inline-block;
	padding:0.5em 0.75em;
	border-radius: 8px;
	box-shadow: 1px 2px 5px #333;
}
p.bloA {
	background: #B6BDE3;
}
p.bloB {
	background: #E3BDB3;
}

.slide-number{
	margin-bottom:10px !important;
	width:100%;
	text-align:center;
	font-size:25px !important;
	background-color:transparent !important;
}
iframe.myclass{
	width:100px;
	height:100px;
	bottom:0;
	left:0;
	position:fixed;
	border:none;
	z-index:99999;
}
h1.raw {
	color: #fff;
	background-image: linear-gradient(90deg,#f35626,#feab3a);
	-webkit-background-clip: text;
	-webkit-text-fill-color: transparent;
	animation: hue 5s infinite linear;
}
@keyframes hue {
	from {
	  filter: hue-rotate(0deg);
	}
	to {
	  filter: hue-rotate(360deg);
	}
}
.progress{
height:14px !important;
}

.progress span{
height:14px !important;
background: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAAMCAIAAAAs6UAAAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyJpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMy1jMDExIDY2LjE0NTY2MSwgMjAxMi8wMi8wNi0xNDo1NjoyNyAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENTNiAoV2luZG93cykiIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6QUNCQzIyREQ0QjdEMTFFMzlEMDM4Qzc3MEY0NzdGMDgiIHhtcE1NOkRvY3VtZW50SUQ9InhtcC5kaWQ6QUNCQzIyREU0QjdEMTFFMzlEMDM4Qzc3MEY0NzdGMDgiPiA8eG1wTU06RGVyaXZlZEZyb20gc3RSZWY6aW5zdGFuY2VJRD0ieG1wLmlpZDpBQ0JDMjJEQjRCN0QxMUUzOUQwMzhDNzcwRjQ3N0YwOCIgc3RSZWY6ZG9jdW1lbnRJRD0ieG1wLmRpZDpBQ0JDMjJEQzRCN0QxMUUzOUQwMzhDNzcwRjQ3N0YwOCIvPiA8L3JkZjpEZXNjcmlwdGlvbj4gPC9yZGY6UkRGPiA8L3g6eG1wbWV0YT4gPD94cGFja2V0IGVuZD0iciI/PovDFgYAAAAmSURBVHjaYvjPwMAAxjMZmBhA9H8INv4P4TPM/A+m04zBNECAAQBCWQv9SUQpVgAAAABJRU5ErkJggg==") repeat-x !important;

}

.progress span:after,
.progress span.nyancat{
	content: "";
	background: url('data:image/gif;base64,R0lGODlhIgAVAKIHAL3/9/+Zmf8zmf/MmZmZmf+Z/wAAAAAAACH/C05FVFNDQVBFMi4wAwEAAAAh/wtYTVAgRGF0YVhNUDw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMy1jMDExIDY2LjE0NTY2MSwgMjAxMi8wMi8wNi0xNDo1NjoyNyAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wTU09Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9tbS8iIHhtbG5zOnN0UmVmPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvc1R5cGUvUmVzb3VyY2VSZWYjIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDpDMkJBNjY5RTU1NEJFMzExOUM4QUM2MDAwNDQzRERBQyIgeG1wTU06RG9jdW1lbnRJRD0ieG1wLmRpZDpCREIzOEIzMzRCN0IxMUUzODhEQjgwOTYzMTgyNTE0QiIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDpCREIzOEIzMjRCN0IxMUUzODhEQjgwOTYzMTgyNTE0QiIgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBQaG90b3Nob3AgQ1M2IChXaW5kb3dzKSI+IDx4bXBNTTpEZXJpdmVkRnJvbSBzdFJlZjppbnN0YW5jZUlEPSJ4bXAuaWlkOkM1QkE2NjlFNTU0QkUzMTE5QzhBQzYwMDA0NDNEREFDIiBzdFJlZjpkb2N1bWVudElEPSJ4bXAuZGlkOkMyQkE2NjlFNTU0QkUzMTE5QzhBQzYwMDA0NDNEREFDIi8+IDwvcmRmOkRlc2NyaXB0aW9uPiA8L3JkZjpSREY+IDwveDp4bXBtZXRhPiA8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') !important;
   width: 34px !important;
   height: 21px !important;
   border: none !important;
   float:right;
   margin-top:-7px;
   margin-right:-10px;
}
</style>

