-- -----------------------------------------------------
-- Schema MEDIA_ROOMS
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `MEDIA_ROOMS` ;
USE `MEDIA_ROOMS` ;

-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`User`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`User` (
  `user_id` INT NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_handle` VARCHAR(255) NULL,
  PRIMARY KEY (`user_id`));


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Host`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Host` (
  `host_id` INT NOT NULL AUTO_INCREMENT,
  `host_username` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`host_id`)
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Room`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Room` (
  `room_id` INT NOT NULL AUTO_INCREMENT,
  `user_owner` INT NOT NULL,
  `host_id` INT NOT NULL,
  `room_code` VARCHAR(255) NOT NULL,
  `start_date` DATETIME NOT NULL,
  `end_date` DATETIME NOT NULL,
  PRIMARY KEY (`room_id`),
  INDEX `fk_Room_User_idx` (`user_owner` ASC) VISIBLE,
  INDEX `fk_Room_Host_idx` (`host_id` ASC) VISIBLE,
  CONSTRAINT `fk_Room_User`
    FOREIGN KEY (`user_owner`)
    REFERENCES `MEDIA_ROOMS`.`User` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Room_Host`
    FOREIGN KEY (`host_id`)
    REFERENCES `MEDIA_ROOMS`.`Host` (`host_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Media`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Media` (
  `media_id` INT NOT NULL AUTO_INCREMENT,
  `url` VARCHAR(255) NOT NULL,
  `title` VARCHAR(200) NOT NULL,
  `artist` VARCHAR(200) NULL,
  `year` INT NULL,
  PRIMARY KEY (`media_id`)
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Tag`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Tag` (
  `tag_id` INT NOT NULL AUTO_INCREMENT,
  `tag` VARCHAR(200) NOT NULL,
  PRIMARY KEY (`tag_id`)
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`MediaTag`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`MediaTag` (
  `media_id` INT NOT NULL,
  `tag_id` INT NOT NULL,
  INDEX `fk_MediaTag_Media_idx` (`media_id` ASC) VISIBLE,
  PRIMARY KEY (`media_id`, `tag_id`),
  INDEX `fk_MediaTag_Tag_idx` (`tag_id` ASC) VISIBLE,
  CONSTRAINT `fk_MediaTag_Media`
    FOREIGN KEY (`media_id`)
    REFERENCES `MEDIA_ROOMS`.`Media` (`media_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_MediaTag_Tag`
    FOREIGN KEY (`tag_id`)
    REFERENCES `MEDIA_ROOMS`.`Tag` (`tag_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`RoomQueue`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`RoomQueue` (
  `room_id` INT NOT NULL,
  `media_id` INT NOT NULL,
  `removed` TINYINT NOT NULL DEFAULT 0,
  `insert_date` DATETIME NOT NULL,
  `is_playing` TINYINT NOT NULL DEFAULT 0,
  INDEX `fk_RoomQueue_Media_idx` (`media_id` ASC) VISIBLE,
  INDEX `fk_RoomQueue_Room_idx` (`room_id` ASC) VISIBLE,
  PRIMARY KEY (`room_id`, `media_id`),
  CONSTRAINT `fk_RoomQueue_Media`
    FOREIGN KEY (`media_id`)
    REFERENCES `MEDIA_ROOMS`.`Media` (`media_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_RoomQueue_Room`
    FOREIGN KEY (`room_id`)
    REFERENCES `MEDIA_ROOMS`.`Room` (`room_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Participant`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Participant` (
  `participant_id` INT NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `participant_handle` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`participant_id`),
  INDEX `fk_Participant_User_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_Participant_User`
    FOREIGN KEY (`user_id`)
    REFERENCES `MEDIA_ROOMS`.`User` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`RoomParticipant`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`RoomParticipant` (
  `participant_id` INT NOT NULL,
  `room_id` INT NOT NULL,
  `join_date` DATETIME NOT NULL,
  INDEX `fk_RoomParticipant_Participant_idx` (`participant_id` ASC) VISIBLE,
  INDEX `fk_RoomParticipant_Room_idx` (`room_id` ASC) VISIBLE,
  PRIMARY KEY (`participant_id`, `room_id`),
  CONSTRAINT `fk_RoomParticipant_Participant`
    FOREIGN KEY (`participant_id`)
    REFERENCES `MEDIA_ROOMS`.`Participant` (`participant_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_RoomParticipant_Room`
    FOREIGN KEY (`room_id`)
    REFERENCES `MEDIA_ROOMS`.`Room` (`room_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Genre`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Genre` (
  `genre_id` INT NOT NULL AUTO_INCREMENT,
  `genre_name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`genre_id`)
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`MediaGenre`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`MediaGenre` (
  `media_id` INT NOT NULL,
  `genre_id` INT NOT NULL,
  INDEX `fk_MediaGenre_Genre_idx` (`genre_id` ASC) VISIBLE,
  INDEX `fk_MediaGenre_Media_idx` (`media_id` ASC) VISIBLE,
  PRIMARY KEY (`genre_id`, `media_id`),
  CONSTRAINT `fk_MediaGenre_Genre`
    FOREIGN KEY (`genre_id`)
    REFERENCES `MEDIA_ROOMS`.`Genre` (`genre_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_MediaGenre_Media`
    FOREIGN KEY (`media_id`)
    REFERENCES `MEDIA_ROOMS`.`Media` (`media_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`MediaOrigin`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`MediaOrigin` (
  `media_id` INT NOT NULL,
  `country` VARCHAR(200) NULL,
  `language` VARCHAR(200) NULL,
  INDEX `fk_MediaOrigin_Media_idx` (`media_id` ASC) VISIBLE,
  PRIMARY KEY (`media_id`),
  CONSTRAINT `fk_MediaOrigin_Media`
    FOREIGN KEY (`media_id`)
    REFERENCES `MEDIA_ROOMS`.`Media` (`media_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);