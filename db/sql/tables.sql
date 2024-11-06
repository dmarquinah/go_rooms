-- -----------------------------------------------------
-- Schema MEDIA_ROOMS
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `MEDIA_ROOMS` ;
USE `MEDIA_ROOMS` ;

-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`User`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`User` (
  `userId` INT NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_handle` VARCHAR(255) NULL,
  PRIMARY KEY (`userId`)
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Host`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Host` (
  `hostId` INT NOT NULL AUTO_INCREMENT,
  `host_username` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`hostId`)
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Room`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Room` (
  `roomId` INT NOT NULL,
  `user_owner` INT NOT NULL,
  `hostId` INT NOT NULL,
  `room_code` VARCHAR(255) NOT NULL,
  `start_date` DATETIME NOT NULL,
  `end_date` DATETIME NOT NULL,
  PRIMARY KEY (`roomId`),
  INDEX `fk_Room_User_idx` (`user_owner` ASC) VISIBLE,
  INDEX `fk_Room_Host_idx` (`hostId` ASC) VISIBLE,
  CONSTRAINT `fk_Room_User`
    FOREIGN KEY (`user_owner`)
    REFERENCES `MEDIA_ROOMS`.`User` (`userId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Room_Host`
    FOREIGN KEY (`hostId`)
    REFERENCES `MEDIA_ROOMS`.`Host` (`hostId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`MediaGenre`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`MediaGenre` (
  `genreId` INT NOT NULL,
  `genre_name` VARCHAR(200) NOT NULL,
  PRIMARY KEY (`genreId`)
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Media`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Media` (
  `mediaId` INT NOT NULL,
  `url` VARCHAR(255) NOT NULL,
  `title` VARCHAR(200) NOT NULL,
  `artist` VARCHAR(200) NULL,
  `genreId` INT NULL,
  PRIMARY KEY (`mediaId`),
  INDEX `fk_Media_MediaGenre_idx` (`genreId` ASC) VISIBLE,
  CONSTRAINT `fk_Media_MediaGenre`
    FOREIGN KEY (`genreId`)
    REFERENCES `MEDIA_ROOMS`.`MediaGenre` (`genreId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`MediaTag`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`MediaTag` (
  `mediaId` INT NOT NULL,
  `tag` VARCHAR(100) NOT NULL,
  INDEX `fk_MediaTag_Media_idx` (`mediaId` ASC) VISIBLE,
  PRIMARY KEY (`mediaId`),
  CONSTRAINT `fk_MediaTag_Media`
    FOREIGN KEY (`mediaId`)
    REFERENCES `MEDIA_ROOMS`.`Media` (`mediaId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`RoomQueue`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`RoomQueue` (
  `roomId` INT NOT NULL,
  `mediaId` INT NOT NULL,
  `removed` TINYINT NOT NULL DEFAULT 0,
  `insert_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_playing` TINYINT NOT NULL DEFAULT 0,
  INDEX `fk_RoomQueue_Media_idx` (`mediaId` ASC) VISIBLE,
  INDEX `fk_RoomQueue_Room_idx` (`roomId` ASC) VISIBLE,
  PRIMARY KEY (`roomId`, `mediaId`),
  CONSTRAINT `fk_RoomQueue_Media`
    FOREIGN KEY (`mediaId`)
    REFERENCES `MEDIA_ROOMS`.`Media` (`mediaId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_RoomQueue_Room`
    FOREIGN KEY (`roomId`)
    REFERENCES `MEDIA_ROOMS`.`Room` (`roomId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`Participant`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`Participant` (
  `participantId` INT NOT NULL,
  `userId` INT NOT NULL,
  `participant_handle` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`participantId`),
  INDEX `fk_Participant_User_idx` (`userId` ASC) VISIBLE,
  CONSTRAINT `fk_Participant_User`
    FOREIGN KEY (`userId`)
    REFERENCES `MEDIA_ROOMS`.`User` (`userId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);


-- -----------------------------------------------------
-- Table `MEDIA_ROOMS`.`RoomParticipant`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `MEDIA_ROOMS`.`RoomParticipant` (
  `participantId` INT NOT NULL,
  `roomId` INT NOT NULL,
  `join_date` DATETIME NOT NULL,
  INDEX `fk_RoomParticipant_Participant_idx` (`participantId` ASC) VISIBLE,
  INDEX `fk_RoomParticipant_Room_idx` (`roomId` ASC) VISIBLE,
  PRIMARY KEY (`participantId`, `roomId`),
  CONSTRAINT `fk_RoomParticipant_Participant`
    FOREIGN KEY (`participantId`)
    REFERENCES `MEDIA_ROOMS`.`Participant` (`participantId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_RoomParticipant_Room`
    FOREIGN KEY (`roomId`)
    REFERENCES `MEDIA_ROOMS`.`Room` (`roomId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);