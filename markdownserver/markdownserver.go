/*
 *
 * Copyright 2023 puzzlemarkdownserver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package markdownserver

import (
	"bytes"
	"context"
	"errors"

	"github.com/dvaumoron/puzzlemarkdownserver/wikilink"
	pb "github.com/dvaumoron/puzzlemarkdownservice"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"go.uber.org/zap"
)

const MarkdownKey = "puzzleMarkdown"

var errInternal = errors.New("internal service error")

// server is used to implement puzzlesaltservice.SaltServer
type server struct {
	pb.UnimplementedMarkdownServer
	md     goldmark.Markdown
	logger *otelzap.Logger
}

func New(logger *otelzap.Logger) pb.MarkdownServer {
	return server{md: goldmark.New(
		goldmark.WithExtensions(extension.GFM, wikilink.Extension),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	), logger: logger}
}

func (s server) Apply(ctx context.Context, request *pb.MarkdownText) (*pb.MarkdownHtml, error) {
	var buf bytes.Buffer
	if err := s.md.Convert([]byte(request.Text), &buf); err != nil {
		s.logger.ErrorContext(ctx, "Failed to transform markdown", zap.Error(err))
		return nil, errInternal
	}
	return &pb.MarkdownHtml{Html: buf.String()}, nil
}
