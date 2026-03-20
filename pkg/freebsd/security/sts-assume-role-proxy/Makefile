PORTNAME=	sts-assume-role-proxy
DISTVERSIONPREFIX=	v
DISTVERSION=	0.1.0
PORTREVISION=	0
CATEGORIES=	security

MAINTAINER=	jesse@mister-webhooks.com
COMMENT=	AWS Security Token Service proxy
WWW=https://github.com/mister-webhooks/sts-assume-role-proxy

LICENSE=	BSD2CLAUSE
LICENSE_FILE=	${WRKSRC}/LICENSE

USES=		go:1.23+,modules

USE_GITHUB=	yes
GH_ACCOUNT=	mister-webhooks

GO_PKGNAME=	github.com/${GH_ACCOUNT}/${PORTNAME}

PLIST_FILES=	bin/${PORTNAME}

.include <bsd.port.mk>
