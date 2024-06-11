CREATE TABLE public.configurations (
    configuration_id serial PRIMARY KEY,
    description character varying(100) NOT NULL,
    engine character varying(50) NOT NULL,
    gearbox_type character varying(50) NOT NULL,
    engine_management_system boolean NOT NULL,
    monitoring_control_system character varying(50) NOT NULL,
    starting_system boolean NOT NULL,
    power_transmission boolean NOT NULL,
    auxiliary_pto character varying(50)[],
    oil_system character varying(50)[],
    fuel_system character varying(50)[],
    cooling_system character varying(50) NOT NULL,
    exhaust_system boolean NOT NULL,
    mounting_system boolean NOT NULL,
    gearbox_option character varying(50)[],
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    last_modified timestamp without time zone DEFAULT now() NOT NULL
);


INSERT INTO configurations(description, engine, gearbox_type, csr_compliance, monitoring_control_system, air_starter, torsionally_resilient_coupling,
    auxiliary_pto,
    oil_system,
    fuel_system,
    cooling_system,
    exhaust_bellows_discharge_rotatable_90,
    resilient_mounts_driving_end,
    gearbox_option,
    created_on,
    last_modified_on
    )
VALUES ('Erste Testkonfiguration', 'test', 'test', true, 'test', false, true, 'test', 'test', 'test', 'test', true, false, 'test', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

SELECT * FROM configurations
ORDER BY configuration_id ASC 


