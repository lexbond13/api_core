BEGIN;

CREATE TABLE club
(
	id bigserial constraint Club_pk primary key,
	name varchar(255) not null,
	tagline varchar(255),
	logo text,
	cover_image text,
	description text not null,
	address varchar(255) not null,
	email varchar(255) not null,
	phone varchar(255) not null,
	rating int default 0,
	status varchar(255) default 'new',
	owner_id bigint constraint Club_user_id_fk references "user" on update cascade on delete restrict,
	created_at timestamp,
	updated_at timestamp
);

create unique index club_name_uindex	on club (name);

COMMIT;
