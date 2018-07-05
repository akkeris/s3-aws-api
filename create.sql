-- Table: public.provision

-- DROP TABLE public.provision;

CREATE TABLE if not exists public.provision
(
  bucketname text NOT NULL,
  location text,
  accesskey_enc text,
  secretkey_enc text,
  billingcode text,
  CONSTRAINT provision_pkey PRIMARY KEY (bucketname)
);

