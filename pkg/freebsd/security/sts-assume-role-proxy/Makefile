PORTNAME=		sts-assume-role-proxy
DISTVERSIONPREFIX=	v
DISTVERSION=		0.1.0
PORTREVISION=		0
CATEGORIES=		security

MAINTAINER=	jesse@mister-webhooks.com
COMMENT=	AWS Security Token Service proxy
WWW=		https://github.com/mister-webhooks/sts-assume-role-proxy

LICENSE=	BSD2CLAUSE
LICENSE_FILE=	${WRKSRC}/LICENSE

USES=		go:1.23+,modules

GO_MODULE=	github.com/mister-webhooks/sts-assume-role-proxy

STS_ASSUME_ROLE_PROXY_USER?=	sts-assume-role-proxy
STS_ASSUME_ROLE_PORXY_GROUP?=	sts-assume-role-proxy
STS_ASSUME_ROLE_PROXY_LOGDIR?=	/var/log/${PORTNAME}

USERS=		${STS_ASSUME_ROLE_PROXY_USER}
GROUPS=		${STS_ASSUME_ROLE_PORXY_GROUP}
USE_RC_SUBR=	sts_assume_role_proxy

PLIST_FILES=	"@sample(${STS_ASSUME_ROLE_PROXY_USER},${STS_ASSUME_ROLE_PORXY_GROUP},0644) etc/sts_assume_role_proxy.conf.sample"\
	bin/${PORTNAME}

SUB_LIST=	GROUP=${STS_ASSUME_ROLE_PORXY_GROUP} \
		USER=${STS_ASSUME_ROLE_PROXY_USER}

post-install:
	${INSTALL_DATA} ${FILESDIR}/sts_assume_role_proxy.conf.sample ${STAGEDIR}${PREFIX}/etc/sts_assume_role_proxy.conf.sample

.include <bsd.port.mk>
