CREATE TABLE IF NOT EXISTS `survey` (
    `id` CHARACTER(16) PRIMARY KEY,
    `title` TEXT NOT NULL,
    `description` TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS `question` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `type` TEXT NOT NULL,
    `statement` TEXT NOT NULL,
    `survey_id` CHARACTER(16) NOT NULL,
    FOREIGN KEY(`survey_id`) REFERENCES `survey`(`id`)
);

CREATE TABLE IF NOT EXISTS `answer` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `answer` TEXT NOT NULL,
    `question_id` INTEGER NOT NULL,
    FOREIGN KEY(`question_id`) REFERENCES `question`(`id`)
);

CREATE TABLE IF NOT EXISTS `response` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `survey_id` CHARACTER(16) NOT NULL,
    `question_id` INTEGER NOT NULL,
    `answer` TEXT NOT NULL,
    FOREIGN KEY(`survey_id`) REFERENCES `survey`(`id`),
    FOREIGN KEY(`question_id`) REFERENCES `question`(`id`)
);
