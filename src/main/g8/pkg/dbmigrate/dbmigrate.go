package dbmigrate

import (
	"context"
	"fmt"
	"os"

	"github.com/Gympass/gcore/v3/gcontext"
	"github.com/Gympass/gcore/v3/glog"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/pkg/errors"

	// postgres import is needed by migrate to connect to database.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// file is necessary to read migrations from filesystem.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	// ErrDirtyMigration is returned when the migrate status is dirty.
	ErrDirtyMigration = errors.New("migration is dirty, needs a manual fix")
)

// Config is used to receive all parameters to apply a migration.
type Config struct {
	Host, Port string
	User, Pass string
	Database   string
	Directory  string
	Logger     glog.Logger
}

// Migrate is used to go up and down with migrations.
type Migrate struct {
	dsn string
	src string

	sd      source.Driver
	migrate *migrate.Migrate
	logger  glog.Logger
}

// New returns a Migrate instance with a connection to a PostgreSQL server.
// The argument opts is an array of options.
// The opts format must by in key=value format.
// Example: sslmode=disable connect_timeout=5
func New(c Config, opts ...string) (*Migrate, error) {
	var query string
	for i := range opts {
		if i == 0 {
			query = opts[0]
			continue
		}

		query = fmt.Sprintf("%s&%s", query, opts[i])
	}

	pgdsn := newDSN(c.User, c.Pass, c.Host, c.Port, c.Database, query)

	srcURL := fmt.Sprintf("file://%s", c.Directory)

	sourceDrv, err := source.Open(srcURL)
	if err != nil {
		return nil, err
	}

	m, err := migrate.New(srcURL, pgdsn)
	if err != nil {
		return nil, err
	}

	return &Migrate{
		dsn:     newDSN("user", "pass", c.Host, c.Port, c.Database, query),
		src:     srcURL,
		sd:      sourceDrv,
		migrate: m,
		logger:  c.Logger,
	}, nil
}

func newDSN(user, pass, host, port, db, query string) string {

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		user,
		pass,
		host,
		port,
		db,
		query,
	)
}

// Up applies all available migrations.
func (m *Migrate) Up() (err error) {
	m.logInfo("Current migration status.")

	_, dirty, err := m.version()
	if err != nil {
		return err
	}

	if dirty {
		if err := m.Down(); err != nil {
			return err
		}
	}

	err = m.migrate.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logInfo("Nothing changed.")
			return nil
		}

		if err := m.Down(); err != nil {
			return err
		}

		// return the m.migrate.Up() error.
		return errors.Wrap(err, "migration up failed")
	}

	m.logInfo("Migration applied.")
	return nil
}

// Down rollback to previous version.
func (m *Migrate) Down() error {
	var (
		mVersion uint
		err      error
	)
	mVersion, _, err = m.version()
	if err != nil {
		return err
	}

	prevVersion, err := m.sd.Prev(mVersion)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if errors.Is(err, os.ErrNotExist) {
		// From now on the code will force dirty to be false for the
		// current version and apply Down() to all migrations.
		// This is an edge case when the first migration fails.

		return m.downAllMigrations(mVersion)
	}

	// Forcing the previous version allows the dirty version to be fixed.
	m.logWarn("Forcing migration.")
	if err := m.migrate.Force(int(prevVersion)); err != nil {
		return errors.Wrap(ErrDirtyMigration, err.Error())
	}

	var dirty bool
	_, dirty, err = m.version()
	if err != nil {
		return err
	}

	if dirty {
		// Going to the previous version did not worked.
		return ErrDirtyMigration
	}

	m.logWarn("Migration Down() applied.")
	return nil
}

// Close is necessary to close the file driver and database connection.
func (m *Migrate) Close() (err error) {
	err = m.sd.Close()

	se, de := m.migrate.Close()
	if se != nil {
		if err != nil {
			err = fmt.Errorf("%s; %s", err, se)
		} else {
			err = se
		}
	}
	if de != nil {
		if err != nil {
			err = fmt.Errorf("%s; %s", err, se)
		} else {
			err = de
		}
	}

	return err
}

// downAllMigrations is used to remove all migrations when the current version
// is the first migration and is dirty.
func (m *Migrate) downAllMigrations(mVersion uint) error {
	first, err := m.sd.First()
	if err != nil {
		return err
	}
	if first != mVersion {
		return errors.Wrapf(
			ErrDirtyMigration,
			"%v is not the first migration (%v) and is missing the previous one",
			mVersion,
			first,
		)
	}
	if err := m.migrate.Force(int(mVersion)); err != nil {
		return errors.Wrap(ErrDirtyMigration, err.Error())
	}

	m.logWarn("Going down with all migrations...")
	if err = m.migrate.Down(); err != nil {
		return err
	}

	m.logWarn("All migrations removed.")

	return nil
}

func (m *Migrate) version() (mVersion uint, dirty bool, err error) {
	mVersion, dirty, err = m.migrate.Version()
	if err != nil {
		if !errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, err
		}
	}

	return mVersion, dirty, nil
}

type logFunc func(context.Context, string, ...any)

func (m *Migrate) logInfo(msg string) {
	m.log(msg, m.logger.Info)
}

func (m *Migrate) logWarn(msg string) {
	m.log(msg, m.logger.Warn)
}

func (m *Migrate) log(msg string, f logFunc) {
	ctx := gcontext.NewContext(context.Background())
	gcontext.AddString(ctx, "migration.data_source_name", m.dsn)
	gcontext.AddString(ctx, "migration.source", m.src)

	version, dirty, err := m.version()
	if err != nil {
		gcontext.AddError(ctx, err)
		m.logger.Error(ctx, "getting migration version")
		return
	}

	gcontext.AddString(ctx, "migration.version", fmt.Sprintf("%v", version))
	gcontext.AddString(ctx, "migration.dirty", fmt.Sprintf("%v", dirty))

	f(ctx, msg)
}
