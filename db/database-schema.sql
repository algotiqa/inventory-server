-- ======================================================================
-- ===   Sql Script for Database : Inventory Server
-- ===
-- === Build : 317
-- ======================================================================

CREATE TABLE trading_session
  (
    id          int           auto_increment,
    username    varchar(32),
    name        varchar(64)   not null,
    session     text          not null,
    created_at  datetime      not null,
    updated_at  datetime,

    primary key(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX trading_sessionIDX1 ON trading_session(username);

-- ======================================================================

CREATE TABLE currency
  (
    id             int           auto_increment,
    code           varchar(16)   unique not null,
    name           varchar(32)   not null,
    symbol         varchar(4)    not null,
    first_date     int           not null,
    last_date      int           not null,
    last_value     double        not null,
    history_ended  tinyint       not null,

    primary key(id)
  )
 ENGINE = InnoDB ;

-- ======================================================================

CREATE TABLE currency_history
  (
    id           int      auto_increment,
    currency_id  int      not null,
    date         int      not null,
    value        double   not null,

    primary key(id),
    unique(currency_id,date),

    foreign key(currency_id) references currency(id)
  )
 ENGINE = InnoDB ;

-- ======================================================================

CREATE TABLE connection
  (
    id                      int           auto_increment,
    username                varchar(32)   not null,
    code                    varchar(8)    not null,
    name                    varchar(32)   not null,
    system_code             varchar(8)    not null,
    system_name             varchar(32)   not null,
    system_config_params    text,
    connected               tinyint       not null,
    supports_data           tinyint       not null,
    supports_broker         tinyint       not null,
    supports_multiple_data  tinyint       not null,
    supports_inventory      tinyint       not null,
    created_at              datetime      not null,
    updated_at              datetime,

    primary key(id),
    unique(username,code)
  )
 ENGINE = InnoDB ;

CREATE INDEX connectionIDX1 ON connection(username);

-- ======================================================================

CREATE TABLE agent_profile
  (
    id              int               auto_increment,
    username        varchar(32)       not null,
    name            varchar(64)       not null,
    host            varchar(255)      not null,
    port            int               not null,
    scan_interval   int               not null,
    scan_folder     varchar(255)      not null,
    file_extension  varchar(16)       not null,
    host_type       varchar(16)       not null,
    ssl_key         varbinary(3072)   not null,
    ssl_cert        varbinary(3072)   not null,
    created_at      datetime          not null,
    updated_at      datetime,

    primary key(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX agent_profileIDX1 ON agent_profile(username);

-- ======================================================================

CREATE TABLE portfolio
  (
    id           int           auto_increment,
    username     varchar(32)   not null,
    name         varchar(64)   not null,
    currency_id  int           not null,
    capital      int           not null,
    target_risk  int           not null,
    created_at   datetime      not null,
    updated_at   datetime,

    primary key(id),

    foreign key(currency_id) references currency(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX portfolioIDX1 ON portfolio(username);

-- ======================================================================

CREATE TABLE outbox
  (
    id         int           auto_increment,
    timestamp  datetime      not null,
    exchange   varchar(64)   not null,
    uuid       varchar(64)   not null,
    payload    mediumblob    not null,
    size       int           not null,

    primary key(id)
  )
 ENGINE = InnoDB ;

-- ======================================================================

CREATE TABLE exchange
  (
    id           int            auto_increment,
    currency_id  int            not null,
    code         varchar(16)    not null,
    name         varchar(64)    not null,
    timezone     varchar(32)    not null,
    url          varchar(255),

    primary key(id),

    foreign key(currency_id) references currency(id)
  )
 ENGINE = InnoDB ;

-- ======================================================================

CREATE TABLE data_product
  (
    id                int           auto_increment,
    exchange_id       int           not null,
    connection_id     int           not null,
    session_id        int           not null,
    username          varchar(32)   not null,
    symbol            varchar(16)   not null,
    name              varchar(64)   not null,
    market_type       char(2)       not null,
    product_type      char(2)       not null,
    months            varchar(16),
    rollover_trigger  varchar(16),
    created_at        datetime      not null,
    updated_at        datetime,

    primary key(id),
    unique(connection_id,username,symbol),

    foreign key(exchange_id) references exchange(id),
    foreign key(connection_id) references connection(id),
    foreign key(session_id) references trading_session(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX data_productIDX1 ON data_product(username);

-- ======================================================================

CREATE TABLE broker_product
  (
    id                  int           auto_increment,
    exchange_id         int           not null,
    connection_id       int           not null,
    username            varchar(32)   not null,
    symbol              varchar(16)   not null,
    name                varchar(64)   not null,
    point_value         double        not null,
    cost_per_operation  double        not null,
    margin_value        double        not null,
    increment           double        not null,
    market_type         char(2)       not null,
    product_type        char(2)       not null,
    created_at          datetime      not null,
    updated_at          datetime,

    primary key(id),
    unique(connection_id,username,symbol),

    foreign key(exchange_id) references exchange(id),
    foreign key(connection_id) references connection(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX broker_productIDX1 ON broker_product(username);

-- ======================================================================

CREATE TABLE trading_system
  (
    id                  int            auto_increment,
    username            varchar(32)    not null,
    data_product_id     int            not null,
    broker_product_id   int            not null,
    trading_session_id  int            not null,
    name                varchar(64)    not null,
    timeframe           int            not null,
    strategy_type       char(2)        not null,
    overnight           tinyint        not null,
    tags                varchar(255),
    agent_profile_id    int,
    external_ref        varchar(64),
    finalized           tinyint        not null,
    in_sample_from      int            not null,
    in_sample_to        int            not null,
    engine_code         varchar(16)    not null,
    created_at          datetime       not null,
    updated_at          datetime,

    primary key(id),

    foreign key(data_product_id) references data_product(id),
    foreign key(broker_product_id) references broker_product(id),
    foreign key(trading_session_id) references trading_session(id),
    foreign key(agent_profile_id) references agent_profile(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX trading_systemIDX1 ON trading_system(username);

-- ======================================================================

