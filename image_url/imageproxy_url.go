package image_url

import (
	"strings"
	"strconv"
)

type imageproxyURL struct {
	imageproxyHost string
	baseURL        string
	bucketName     string
	omitBaseURL    bool
}

func (ip *imageproxyURL) Generate(hashValue string, opt ...Option) string {
	opts := defaultURLOptions
	for _, o := range opt {
		o(&opts)
	}
	sb := strings.Builder{}
	sb.WriteString(ip.imageproxyHost)
	sb.WriteRune('/')
	sb.WriteString(ip.buildOptionsStr(&opts))
	sb.WriteRune('/')
	if !ip.omitBaseURL {
		sb.WriteString(ip.baseURL)
		sb.WriteRune('/')
	}

	sb.WriteString(ip.bucketName)
	sb.WriteRune('/')
	sb.WriteString(hashValue)
	return sb.String()
}

func (ip *imageproxyURL) buildOptionsStr(opts *options) string {
	opt := strings.Builder{}
	if opts.width != 0 {
		opt.WriteString(strconv.Itoa(int(opts.width)))
		opt.WriteRune('x')
	} else if opts.widthPercent != 0 {
		opt.WriteString(strconv.FormatFloat(float64(opts.widthPercent), 'f', -1, 32))
		opt.WriteRune('x')
	}

	if opts.height != 0 {
		if opt.Len() <= 0 {
			opt.WriteRune('x')
		}
		opt.WriteString(strconv.Itoa(int(opts.height)))
	} else if opts.heightPercent != 0 {
		if opt.Len() <= 0 {
			opt.WriteRune('x')
		}
		opt.WriteString(strconv.FormatFloat(float64(opts.heightPercent), 'f', -1, 32))
	}

	if opts.quality != 100 {
		if opt.Len() > 0 {
			opt.WriteRune(',')
		}
		opt.WriteRune('q')
		opt.WriteString(strconv.Itoa(int(opts.quality)))
	}
	return opt.String()
}

func NewImageproxyURL(imageproxyHost, baseURL, bucketName string, omitBaseURL bool) URL {
	return &imageproxyURL{
		imageproxyHost: strings.TrimRight(imageproxyHost, "/"),
		baseURL:        strings.TrimRight(baseURL, "/"),
		bucketName:     bucketName,
		omitBaseURL:    omitBaseURL,
	}
}
