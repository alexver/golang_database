package compute

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alexver/golang_database/internal/compute/analyzer"
	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestCreateComputeLayer(t *testing.T) {
	type args struct {
		parser parser.ParserInterface
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "check analyzers map initialization",
			args: args{parser: nil, logger: zaptest.NewLogger(t)},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateComputeLayer(tt.args.parser, tt.args.logger); !reflect.DeepEqual(len(got.GetAnalyzers()), tt.want) {
				t.Errorf("CreateComputeLayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

// parser mock
type parserMock struct {
	mock.Mock
}

func (m *parserMock) ParseStringToQuery(command string) parser.Query {
	args := m.Called(command)
	return args.Get(0).(parser.Query)
}

// analyzer mock
type analyzerMock struct {
	mock.Mock
}

func (m *analyzerMock) Name() string {
	args := m.Called()
	return args.String(0)
}
func (m *analyzerMock) Description() string {
	args := m.Called()
	return args.String(0)
}
func (m *analyzerMock) Usage() string {
	args := m.Called()
	return args.String(0)
}
func (m *analyzerMock) Supports(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}
func (m *analyzerMock) Validate(query parser.Query) error {
	args := m.Called(query)
	return args.Error(0)
}
func (m *analyzerMock) Run(query parser.Query) (any, error) {
	args := m.Called(query)
	return args.String(0), args.Error(1)
}

func TestCompute_HandleQuery(t *testing.T) {
	type fields struct {
		parser    parser.ParserInterface
		analyzers map[string]analyzer.AnalyzerInterface
		logger    *zap.Logger
	}
	type args struct {
		queryStr string
	}

	list := map[string]analyzer.AnalyzerInterface{}
	list["EXIT"] = analyzer.NewExit(zaptest.NewLogger(t))

	parserMock1 := new(parserMock)
	parserMock1.On("ParseStringToQuery", "EXIT").Once().Return(parser.CreateQuery("EXIT", []string{}))

	parserMock2 := new(parserMock)
	parserMock2.On("ParseStringToQuery", "TEST").Once().Return(parser.CreateQuery("TEST", []string{}))

	query3 := parser.CreateQuery("MOCK", []string{})
	parserMock3 := new(parserMock)
	parserMock3.On("ParseStringToQuery", "MOCK").Return(query3)

	analyzerMock1 := new(analyzerMock)
	analyzerMock1.On("Supports", mock.Anything).Once().Return(true)
	analyzerMock1.On("Validate", query3).Once().Return(errors.New("wrong param"))
	list1 := map[string]analyzer.AnalyzerInterface{}
	list1["MOCK"] = analyzerMock1

	analyzerMock2 := new(analyzerMock)
	analyzerMock2.On("Supports", mock.Anything).Once().Return(true)
	analyzerMock2.On("Validate", mock.Anything).Once().Return(nil)
	analyzerMock2.On("Run", query3).Once().Return("ok", nil)
	analyzerMock2.On("Name").Once().Return("MOCK")
	list2 := map[string]analyzer.AnalyzerInterface{}
	list2["MOCK"] = analyzerMock2

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      any
		wantErr   bool
		errString string
	}{
		{
			name:      "empty analyzer list",
			fields:    fields{parser: parserMock1, analyzers: map[string]analyzer.AnalyzerInterface{}, logger: zaptest.NewLogger(t)},
			args:      args{queryStr: "EXIT"},
			want:      nil,
			wantErr:   true,
			errString: "command EXIT is unknown",
		},
		{
			name:      "no analyzers to process query",
			fields:    fields{parser: parserMock2, analyzers: list, logger: zaptest.NewLogger(t)},
			args:      args{queryStr: "TEST"},
			want:      nil,
			wantErr:   true,
			errString: "command TEST is unknown",
		},
		{
			name:      "validation failed",
			fields:    fields{parser: parserMock3, analyzers: list1, logger: zaptest.NewLogger(t)},
			args:      args{queryStr: "MOCK"},
			want:      nil,
			wantErr:   true,
			errString: "command MOCK arguments are invalid: wrong param",
		},
		{
			name:      "command processed",
			fields:    fields{parser: parserMock3, analyzers: list2, logger: zaptest.NewLogger(t)},
			args:      args{queryStr: "MOCK"},
			want:      "ok",
			wantErr:   false,
			errString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compute{
				parser:    tt.fields.parser,
				analyzers: tt.fields.analyzers,
				logger:    tt.fields.logger,
			}
			got, err := c.HandleQuery(tt.args.queryStr)
			if (err != nil) != tt.wantErr || (err != nil && err.Error() != tt.errString) {
				t.Errorf("Compute.HandleQuery() error = %v, wantErr %v, error message: %s", err, tt.wantErr, tt.errString)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compute.HandleQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompute_RegisterAnalyzer(t *testing.T) {
	type fields struct {
		parser    parser.ParserInterface
		analyzers map[string]analyzer.AnalyzerInterface
		logger    *zap.Logger
	}
	type args struct {
		analyzer analyzer.AnalyzerInterface
	}

	list := map[string]analyzer.AnalyzerInterface{}
	list["EXIT"] = analyzer.NewExit(zaptest.NewLogger(t))

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		errMessage string
		wantCount  int
	}{
		{
			name:       "empty analyzer",
			fields:     fields{parser: nil, analyzers: map[string]analyzer.AnalyzerInterface{}, logger: zaptest.NewLogger(t)},
			args:       args{analyzer: nil},
			wantErr:    true,
			errMessage: "analyzer is not defined",
			wantCount:  0,
		},
		{
			name:       "analyzer already exist",
			fields:     fields{parser: nil, analyzers: list, logger: zaptest.NewLogger(t)},
			args:       args{analyzer: analyzer.NewExit(zaptest.NewLogger(t))},
			wantErr:    true,
			errMessage: "analyzer 'EXIT' already registered",
			wantCount:  1,
		},
		{
			name:       "add analyzer",
			fields:     fields{parser: nil, analyzers: map[string]analyzer.AnalyzerInterface{}, logger: zaptest.NewLogger(t)},
			args:       args{analyzer: analyzer.NewExit(zaptest.NewLogger(t))},
			wantErr:    false,
			errMessage: "",
			wantCount:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compute{
				parser:    tt.fields.parser,
				analyzers: tt.fields.analyzers,
				logger:    tt.fields.logger,
			}
			err := c.RegisterAnalyzer(tt.args.analyzer)
			if (err != nil) != tt.wantErr || (err != nil && err.Error() != tt.errMessage) {
				t.Errorf("Compute.RegisterAnalyzer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(c.analyzers) != tt.wantCount {
				t.Errorf("Compute.RegisterAnalyzer() analyze count = %d, wantCount %d", len(c.analyzers), tt.wantCount)
			}
		})
	}
}

func TestCompute_GetAnalyzers(t *testing.T) {
	type fields struct {
		parser    parser.ParserInterface
		analyzers map[string]analyzer.AnalyzerInterface
		logger    *zap.Logger
	}
	list := map[string]analyzer.AnalyzerInterface{}
	list["EXIT"] = analyzer.NewExit(zaptest.NewLogger(t))

	tests := []struct {
		name      string
		fields    fields
		wantCount int
	}{
		{
			name:      "empty list",
			fields:    fields{parser: nil, analyzers: map[string]analyzer.AnalyzerInterface{}, logger: zaptest.NewLogger(t)},
			wantCount: 0,
		},
		{
			name:      "prefilled list",
			fields:    fields{parser: nil, analyzers: list, logger: zaptest.NewLogger(t)},
			wantCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compute{
				parser:    tt.fields.parser,
				analyzers: tt.fields.analyzers,
				logger:    tt.fields.logger,
			}
			if got := c.GetAnalyzers(); !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("Compute.GetAnalyzers() = %v (count: %d), want count %d", got, len(got), tt.wantCount)
			}
		})
	}
}
