---
title: Don't bieleve "Believe me" - the pitfalls of authentication
tags: Slide, Certificates, Signature, Authentication, Authorization 
description: View the slide with "Slide Mode".
slideOptions:
  spotlight:
    enabled: false
  allottedMinutes: 40
---
<small>不能相信的相信我之術: 網路連線雙方驗證時，那些意想不到的坑</small>
===

<!-- .slide: data-background-color="pink" -->
<!-- .slide: data-transition="zoom" -->

Don't bieleve "Believe me"
the pitfalls of authentication

> [name=郭學聰 Hsueh-Tsung Kuo]
> [time=Sat, 15 Oct 2022] [color=red]

###### CC BY-SA 4.0

---

<!-- .slide: data-transition="convex" -->

## Who am I?

![fieliapm](https://www.gravatar.com/avatar/2aef78f04240a6ac9ccd473ba1cbd1e3?size=2048 =384x384)

<small>Someone (who?) said:
a game programmer should be able to draw cute anime character(?)</small>

----

<!-- .slide: data-transition="convex" -->

* A programmer from game company in Taiwan.
* Backend (and temporary frontend) engineer.
* Usually develop something related to my work in Python, Ruby, ECMAScript, Golang, C#.
* ECMAScript hater since **Netscape** is dead.
* Built CDN-aware game asset update system.
* Business large passenger vehicle driver. :bus: :racing_motorcycle:
* Draw cute anime character in spare time.

---

<!-- .slide: data-transition="convex" -->

## Outline

----

<!-- .slide: data-transition="convex" -->

4. Introduction
    * Encrypted data exchanging
    * Trusted data exchanging
5. Trusted data exchanging
    * User
      * Password
      * One-time password
      * MFA
      * OAuth 1/2
      * Signature
        * HMAC
        * RSA

----

<!-- .slide: data-transition="convex" -->

5. Trusted data exchanging
    * Service
      * SSL
      * TLS
      * Public key certificate
        * RSA
        * ECC: ECDSA/EdDSA/Ed25519
      * CA
      * Chain of trust
    * Peer to peer
      * Centralized
      * Distributed

----

<!-- .slide: data-transition="convex" -->

6. When trust is broken...
    * User
      * Online account recovering
      * Offline
    * Server
      * Cert revoke
        * CRL
        * OCSP &amp; OCSP Stapling
      * The real world - revoking might not work
        * DV/OV/EV
        * Client: Chrome/Edge,Firefox,Safari,IE,Android

----

<!-- .slide: data-transition="convex" -->

7. Conclusion
8. Resource
9. Q&A

---

<!-- .slide: data-transition="convex" -->

## Introduction

----

<!-- .slide: data-transition="convex" -->

### Encrypted data exchanging

* To ensure that only connection peer can understand the content.
  * Man in the middle could knows they exchange something.
  * Contents might be guessed.

----

<!-- .slide: data-transition="convex" -->

### Trusted data exchanging

* To ensure that the connection target is really the target.

----

<!-- .slide: data-transition="convex" -->

### Trusted data exchanging

:hash: "How to know that connection target is really the target?"

---

<!-- .slide: data-transition="convex" -->

## Trusted data exchanging

* User
* Service
* Peer to peer

----

<!-- .slide: data-transition="convex" -->

### User

* Password
* One-time password
* MFA
* OAuth 1/2
* Signature

----

<!-- .slide: data-transition="convex" -->

#### Password

* A secret data used to confirm a user's identity.
* Typically a string of characters.
* Example:
  * "qawsedrftgyhujikolp" <!-- .element: class="fragment" data-fragment-index="1" --> ![rinchan](https://stickershop.line-scdn.net/stickershop/v1/sticker/31272487/android/sticker.png)

----

<!-- .slide: data-transition="convex" -->

#### One-time password

* A password that is valid for only one login session or transaction.
* Example:
  * SMS OTP
  * e-mail OTP
  * Security token for online game

----

<!-- .slide: data-transition="convex" -->

##### Security token

![sqex_security_token](https://www.square-enix-games.com/documents/static/images/otp_en/img03.jpg)
<small>https://www.square-enix-games.com/en_US/seaccount/otp</small>

----

<!-- .slide: data-transition="convex" -->

#### MFA

* Multi-factor authentication
  * Subset: Two-factor authentication.
  * Successfully presenting two or more pieces of evidence.
    * Ex: password + one-time code

----

<!-- .slide: data-transition="convex" -->

#### OAuth 1/2

* "Open Authorization"
* Grant websites or apps access to personal information on other websites without giving private evidence.

![oauth2](https://oauth.net/images/oauth-2-sm.png)
<small>https://oauth.net/2/</small>

----

<!-- .slide: data-transition="convex" -->

#### Signature

![digital_signature](https://upload.wikimedia.org/wikipedia/commons/thumb/7/78/Private_key_signing.svg/491px-Private_key_signing.svg.png)
<small>https://en.wikipedia.org/wiki/Digital_signature</small>

----

<!-- .slide: data-transition="convex" -->

#### Signature

* Verify the authenticity of digital data by asymmetric cryptography.
  * Signer
    * Calculate hash of digital data.
    * Encrypt hash using private key.
  * Verifier
    * Calculate hash of digital data.
    * Verify hash using public key.

----

<!-- .slide: data-transition="convex" -->

#### Signature

* Only specific key owner can produce the content.
* Own the specific key indicates who it is.
* Suitable for bot :arrow_right: service.
* Example:
  * HMAC
    * Amazon S3
  * RSA
    * Google Cloud Storage

----

<!-- .slide: data-transition="convex" -->

#### Signature

* HMAC
    * Peers share the same key.
* RSA
    * Service provider owns public key.
    * Service user owns private key.

----

<!-- .slide: data-transition="convex" -->

### Service

* SSL/TLS
  * Handshake by asymmetric cryptography.
    * Tell client it is genuine and integrity.
      * By **public key certificate**.
    * Exchange symmetric key.
  * Transfer data by symmetric cryptography.

----

<!-- .slide: data-transition="convex" -->

#### SSL

* Secure Sockets Layer
* Introduced by Netscape.
* SSL 1.0 2.0 3.0

----

<!-- .slide: data-transition="convex" -->

#### TLS

* Transport Layer Security
* IETF refines from SSL 3.0 and standardizes it.
* TLS 1.0 1.1 1.2 1.3

----

<!-- .slide: data-transition="convex" -->

#### Public key certificate

* Identity document with signature.
* To prove it and its key are genuine. 
* Issue
  * Self-signed
  * Sign by CA

----

<!-- .slide: data-transition="convex" -->

#### Public key certificate

* Format
  * X.509
    * ITU standard
    * The format of public key certificates.
* Signature Algorithm
  * ~~DSA~~ (too weak)
  * RSA
  * ECDSA/EdDSA/Ed25519

----

<!-- .slide: data-transition="convex" -->

##### RSA

* Based on modulus &amp; factorization.
* Encryption
  * Encrypt using public key.
  * Decrypt using private key.
* Signature
  * Sign using private key.
  * Verify using public key.

----

<!-- .slide: data-transition="convex" -->

##### ECC (Elliptic Curve Cryptography)

* Based on modulus &amp; elliptic curve.
* Less calculation &amp; memory cost.
* Harder to guess.
* Implementation: ECDSA/EdDSA/Ed25519

----

<!-- .slide: data-transition="convex" -->

![elliptic_curve_cryptography](https://www.allaboutcircuits.com/uploads/articles/Curve_Cryptography_fig03.gif)
<small>https://www.allaboutcircuits.com/technical-articles/elliptic-curve-cryptography-in-embedded-systems/</small>

----

<!-- .slide: data-transition="convex" -->

#### CA

* Certificate Authority
* Sign certificate for websites.
* The certificate which represents CA is self-signed.
  * Distribute certificates with browsers.
    * Browser should be downloaded from trusted origins.

----

<!-- .slide: data-transition="convex" -->

![certificate](https://upload.wikimedia.org/wikipedia/commons/thumb/6/65/PublicKeyCertificateDiagram_It.svg/1018px-PublicKeyCertificateDiagram_It.svg.png =640x453)
<small>https://en.wikipedia.org/wiki/Certificate_authority</small>

----

<!-- .slide: data-transition="convex" -->

#### Chain of trust

* The root certificate tell that ends are genuine.
* Clients validate public key from the end entity up to the root certificate.

----

<!-- .slide: data-transition="convex" -->

![chain_of_trust](https://upload.wikimedia.org/wikipedia/commons/thumb/0/02/Chain_Of_Trust.svg/640px-Chain_Of_Trust.svg.png)
<small>https://en.wikipedia.org/wiki/Chain_of_trust</small>

----

<!-- .slide: data-transition="convex" -->

### Peer to peer

* Centralized
* Distributed

----

<!-- .slide: data-transition="convex" -->

#### Centralized

* Certificate signed by CA
* Example:
  * VPN using certificates signed by CA

----

<!-- .slide: data-transition="convex" -->

#### Distributed

* Pre-shared key or certificate
  * Offline acquired
  * Online acquired
    * Based on trusted and secured connections.
* Example:
  * Various types of VPN protocols
  * Blockchain (cryptocurrencies, NFT)
  * BitTorrent(!?)

---

<!-- .slide: data-transition="convex" -->

## When trust is broken...

Case study:
<small>https://techcrunch.com/2021/11/15/hpe-aruba-data-breach/</small>

----

<!-- .slide: data-transition="convex" -->

### User

* Online account recovering
* Offline

----

<!-- .slide: data-transition="convex" -->

#### Online account recovering

* Detected
  * Temporarily deactivate account.
    * Then wait for account restoring.
* Stolen
  * Call customer service.
    * The account might be stolen forever.

----

<!-- .slide: data-transition="convex" -->

#### Offline

* Go to the counter and show your ID card.

----

<!-- .slide: data-transition="convex" -->

### Server

* Cert expired
* Cert revoke

----

<!-- .slide: data-transition="convex" -->

#### Cert revoke

* CRL
* OCSP
* OCSP Stapling

----

<!-- .slide: data-transition="convex" -->

##### CRL

* Certificate Revocation List
* A black list of certificates published from CA.

----

<!-- .slide: data-transition="convex" -->

##### OCSP

* Online Certificate Status Protocol
* Send request to obtain the revocation status of an X.509 certificate.
* Impact connection performance.
* Leak privacy (others know who accessed something).

----

<!-- .slide: data-transition="convex" -->

##### OCSP Stapling

* Cache OCSP result and integrate it in TLS handshake.

----

<!-- .slide: data-transition="convex" -->

```sequence
client->server: client hello
note right of server: valid OCSP response?
note right of server: Y -> use cached OCSP response
server->CA: N -> OCSP request
CA->server: OCSP response

server->client: server hello
note left of server: certificate/certificate status
note left of client: valid certificate?
client->server: Y -> complete handshake
client->server: N -> abort handshake
```

----

<!-- .slide: data-transition="convex" -->

#### The real world - revoking might not work

* Browsers might not ask revocation status from CRL/OCSP.
  * **Browsers maintain their own certificate revoke list or behavior**. <!-- .element: class="fragment" data-fragment-index="1" -->
    * To initialize connection faster. <!-- .element: class="fragment" data-fragment-index="2" -->

----

<!-- .slide: data-transition="convex" -->

##### DV/OV/EV

* Domain Validation
  * Control the domain
* Organization Validation
  * Registered organization
* Extended Validation
  * Business detail (accountant, etc)

:moneybag: :arrow_right: VIP <!-- .element: class="fragment" data-fragment-index="1" -->

----

<!-- .slide: data-transition="convex" -->

##### Client

<small>

| Browser | Cert DV | Cert EV | Comment |
| ------- | ------- | ------- | ------- |
| Chrome/Edge (Windows) | :x: | :x: (?) | always ignored(?) |
| Firefox (OCSP on (default)) | :o: | :o: | OneCRL |
| Firefox (OCSP off) | :x: | :x: | OneCRL |
| Safari | :o: | :o: | :thinking_face: |
| Chrome/Edge (MacOS) | :o: | :o: | based on MacOS |
| IE | :o: | :o: | very strict | |
| Android | :x: | :x: | always ignored |

</small>

---

<!-- .slide: data-transition="convex" -->

## Conclusion

----

<!-- .slide: data-transition="convex" -->

### Experience

* Authentication is not as strict as you think.
* When handling authentication, doubt it first.

----

<!-- .slide: data-transition="convex" -->

### Slogan

:hash: {質疑是資安的根本|<big>Doubt is essential to information security</big>}

> [name=郭學聰 Hsueh-Tsung Kuo] [time=2022_10_15] [color=red] :notebook:

---

<!-- .slide: data-transition="convex" -->

## Resource 

----

<!-- .slide: data-transition="convex" -->

### Reference

* How Do Browsers Handle Revoked SSL/TLS Certificates?
  * <small>https://www.ssl.com/blogs/how-do-browsers-handle-revoked-ssl-tls-certificates/</small>
* No, don't enable revocation checking (19 Apr 2014)
  * <small>https://www.imperialviolet.org/2014/04/19/revchecking.html</small>
* OCSP &amp; CRL 介紹
  * <small>http://ijecorp.blogspot.com/2016/01/ocsp-crl.html</small>

---

<!-- .slide: data-transition="zoom" -->

## Q&A

---

<style>
.reveal {
    background: #FFDFEF;
    color: black;
}
.reveal h2,
.reveal h3,
.reveal h4,
.reveal h5,
.reveal h6 {
    color: black;
}
.reveal code {
    font-size: 18px !important;
    line-height: 1.2;
}

.progress div{
height:14px !important;
background: hotpink !important;
}

// original template

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

/*.slide-number{
	margin-bottom:10px !important;
	width:100%;
	text-align:center;
	font-size:25px !important;
	background-color:transparent !important;
}*/
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

