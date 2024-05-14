package compute

import (
	"fmt"

	"github.com/alexver/golang_database/internal/compute/analyzer"
	"github.com/alexver/golang_database/internal/compute/parser"
	"go.uber.org/zap"
)

type ComputeInterface interface {
	HandleQuery(queryStr string) (*parser.Query, error)
	RegisterAnalyzer(analyzer analyzer.AnalyzerInterface) error
	GetAnalyzers() []analyzer.AnalyzerInterface
}

type Compute struct {
	parser    parser.ParserInterface
	analyzers map[string]analyzer.AnalyzerInterface

	logger *zap.Logger
}

func CreateComputeLayer(parser parser.ParserInterface, logger *zap.Logger) ComputeInterface {
	return &Compute{
		parser:    parser,
		logger:    logger,
		analyzers: make(map[string]analyzer.AnalyzerInterface),
	}
}

func (c *Compute) HandleQuery(queryStr string) (*parser.Query, error) {
	query := c.parser.ParseStringToQuery(queryStr)

	for _, analyzer := range c.analyzers {
		if !analyzer.Supports(query.GetCommand()) {
			continue
		}
		if validateErr := analyzer.Validate(query); validateErr != nil {
			c.logger.Error(fmt.Sprintf("Command '%s' argments validation error: %s", query.GetCommand(), validateErr.Error()))

			return nil, fmt.Errorf("command %s arguments are invalid: %s", query.GetCommand(), validateErr.Error())
		}

		c.logger.Debug(fmt.Sprintf("Command %s is processing by %s analzer", query.GetCommand(), analyzer.Name()))

		query = analyzer.NormalizeQuery(query)

		return query, nil
	}

	c.logger.Error(fmt.Sprintf("Command %s is unknown, didn't find any analyzer to process", query.GetCommand()))

	return nil, fmt.Errorf("command %s is unknown", query.GetCommand())
}

func (c *Compute) RegisterAnalyzer(analyzer analyzer.AnalyzerInterface) error {
	if analyzer == nil {
		c.logger.Error("Register Analyzer error: analyzer is not defined")

		return fmt.Errorf("analyzer is not defined")
	}

	_, ok := c.analyzers[analyzer.Name()]
	if ok {
		c.logger.Error(fmt.Sprintf("Register Analyzer error: '%s' analyzer already registered", analyzer.Name()))

		return fmt.Errorf("analyzer '%s' already registered", analyzer.Name())
	}

	c.analyzers[analyzer.Name()] = analyzer

	return nil
}

func (c *Compute) GetAnalyzers() []analyzer.AnalyzerInterface {
	result := []analyzer.AnalyzerInterface{}

	for _, analyzer := range c.analyzers {
		result = append(result, analyzer)
	}

	return result
}
