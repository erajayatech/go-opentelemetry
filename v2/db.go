package gootel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	DBSystemKey           = attribute.Key("db.system")
	DBNameKey             = attribute.Key("db.name")
	DBStatementKey        = attribute.Key("db.statement")
	DBOperationKey        = attribute.Key("db.operation")
	DBUserKey             = attribute.Key("db.user")
	DBConnectionStringKey = attribute.Key("db.connection_string")
	DBDurationKey         = attribute.Key("db.duration_ms")
	DBSuccessKey          = attribute.Key("db.success")
	DBErrorKey            = attribute.Key("db.error")
	DBRowCountKey         = attribute.Key("db.rows_affected")
)

// TraceDBOperation creates a span for database operations
func TraceDBOperation(ctx context.Context, system, name, operation, statement string) (context.Context, trace.Span) {
	tracer := otel.Tracer("database")

	attrs := []attribute.KeyValue{
		DBSystemKey.String(system),
		DBNameKey.String(name),
		DBOperationKey.String(operation),
		DBStatementKey.String(statement),
		SpanKindKey.String("client"),
	}

	return tracer.Start(ctx, "DB "+operation, trace.WithAttributes(attrs...))
}

// RecordDBSuccess records successful database operation with duration
func RecordDBSuccess(span trace.Span, rowsAffected int64, duration time.Duration) {
	span.SetAttributes(
		DBSuccessKey.Bool(true),
		DBDurationKey.Int64(duration.Milliseconds()),
		DBRowCountKey.Int64(rowsAffected),
	)
	span.SetStatus(codes.Ok, "Database operation succeeded")
}

// RecordDBError records failed database operation with error details
func RecordDBError(span trace.Span, err error) {
	span.SetAttributes(
		DBSuccessKey.Bool(false),
		DBErrorKey.String(err.Error()),
	)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// RecordDBQueryStats records additional database query statistics
func RecordDBQueryStats(span trace.Span, queryTime time.Duration, rowsScanned int64) {
	span.SetAttributes(
		attribute.Key("db.query_time_ms").Int64(queryTime.Milliseconds()),
		attribute.Key("db.rows_scanned").Int64(rowsScanned),
	)
}

// RecordDBQueryEvent logs the query as an event in the span for debugging purposes
// This makes the query visible in tracing UI as a span event
func RecordDBQueryEvent(span trace.Span, statement string) {
	attrs := []attribute.KeyValue{
		DBStatementKey.String(statement),
	}
	span.AddEvent("db.query", trace.WithAttributes(attrs...))
}

// SetupGORMTracing initializes GORM callbacks for enhanced database tracing
// dbSystem should be one of: "postgresql", "mysql", "sqlite", "sqlserver", "mongodb", "couchbase", "cassandra", etc.
func SetupGORMTracing(db *gorm.DB, dbSystem, dbName string) {
	callback := db.Callback()

	callback.Query().Before("gorm:query").Register("otel:start_query", beforeQuery(dbSystem, dbName))
	callback.Query().After("gorm:query").Register("otel:end_query", afterQuery)

	callback.Create().Before("gorm:create").Register("otel:start_create", beforeCreate(dbSystem, dbName))
	callback.Create().After("gorm:create").Register("otel:end_create", afterCreate)

	callback.Update().Before("gorm:update").Register("otel:start_update", beforeUpdate(dbSystem, dbName))
	callback.Update().After("gorm:update").Register("otel:end_update", afterUpdate)

	callback.Delete().Before("gorm:delete").Register("otel:start_delete", beforeDelete(dbSystem, dbName))
	callback.Delete().After("gorm:delete").Register("otel:end_delete", afterDelete)
}

func beforeQuery(dbSystem, dbName string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		ctx, span := TraceDBOperation(db.Statement.Context, dbSystem, dbName, "SELECT", db.Statement.SQL.String())
		db.InstanceSet("otel_span", span)
		db.InstanceSet("otel_start_time", time.Now())
		db.Statement.Context = ctx
	}
}

func afterQuery(db *gorm.DB) {
	if span, ok := db.InstanceGet("otel_span"); ok {
		if s, ok := span.(trace.Span); ok {
			defer s.End()

			if db.Statement.SQL.String() != "" {
				RecordDBQueryEvent(s, db.Statement.SQL.String())
			}

			if startTime, ok := db.InstanceGet("otel_start_time"); ok {
				if st, ok := startTime.(time.Time); ok {
					duration := time.Since(st)

					if db.Error != nil {
						RecordDBError(s, db.Error)
					} else {
						rowsAffected := db.Statement.RowsAffected
						RecordDBSuccess(s, rowsAffected, duration)
					}
				}
			}
		}
	}
}

func beforeCreate(dbSystem, dbName string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		ctx, span := TraceDBOperation(db.Statement.Context, dbSystem, dbName, "INSERT", db.Statement.SQL.String())
		db.InstanceSet("otel_span", span)
		db.InstanceSet("otel_start_time", time.Now())
		db.Statement.Context = ctx
	}
}

func afterCreate(db *gorm.DB) {
	if span, ok := db.InstanceGet("otel_span"); ok {
		if s, ok := span.(trace.Span); ok {
			defer s.End()

			if db.Statement.SQL.String() != "" {
				RecordDBQueryEvent(s, db.Statement.SQL.String())
			}

			if startTime, ok := db.InstanceGet("otel_start_time"); ok {
				if st, ok := startTime.(time.Time); ok {
					duration := time.Since(st)

					if db.Error != nil {
						RecordDBError(s, db.Error)
					} else {
						rowsAffected := db.Statement.RowsAffected
						RecordDBSuccess(s, rowsAffected, duration)
					}
				}
			}
		}
	}
}

func beforeUpdate(dbSystem, dbName string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		ctx, span := TraceDBOperation(db.Statement.Context, dbSystem, dbName, "UPDATE", db.Statement.SQL.String())
		db.InstanceSet("otel_span", span)
		db.InstanceSet("otel_start_time", time.Now())
		db.Statement.Context = ctx
	}
}

func afterUpdate(db *gorm.DB) {
	if span, ok := db.InstanceGet("otel_span"); ok {
		if s, ok := span.(trace.Span); ok {
			defer s.End()

			if db.Statement.SQL.String() != "" {
				RecordDBQueryEvent(s, db.Statement.SQL.String())
			}

			if startTime, ok := db.InstanceGet("otel_start_time"); ok {
				if st, ok := startTime.(time.Time); ok {
					duration := time.Since(st)

					if db.Error != nil {
						RecordDBError(s, db.Error)
					} else {
						rowsAffected := db.Statement.RowsAffected
						RecordDBSuccess(s, rowsAffected, duration)
					}
				}
			}
		}
	}
}

func beforeDelete(dbSystem, dbName string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		ctx, span := TraceDBOperation(db.Statement.Context, dbSystem, dbName, "DELETE", db.Statement.SQL.String())
		db.InstanceSet("otel_span", span)
		db.InstanceSet("otel_start_time", time.Now())
		db.Statement.Context = ctx
	}
}

func afterDelete(db *gorm.DB) {
	if span, ok := db.InstanceGet("otel_span"); ok {
		if s, ok := span.(trace.Span); ok {
			defer s.End()

			if db.Statement.SQL.String() != "" {
				RecordDBQueryEvent(s, db.Statement.SQL.String())
			}

			if startTime, ok := db.InstanceGet("otel_start_time"); ok {
				if st, ok := startTime.(time.Time); ok {
					duration := time.Since(st)

					if db.Error != nil {
						RecordDBError(s, db.Error)
					} else {
						rowsAffected := db.Statement.RowsAffected
						RecordDBSuccess(s, rowsAffected, duration)
					}
				}
			}
		}
	}
}
