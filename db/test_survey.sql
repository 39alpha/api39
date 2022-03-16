INSERT INTO `survey`(`id`, `title`, `description`) VALUES (
    'f4a92aeee9f1d000',
    'A title',
    'This is your survey. There are many like it, but this one is yours.'
);

INSERT INTO `question`(`type`, `statement`, `survey_id`) VALUES
('short-answer', 'What is your name?', 'f4a92aeee9f1d000'),
('email', 'What is your email address?', 'f4a92aeee9f1d000'),
('multiple-choice', 'How old are you?', 'f4a92aeee9f1d000'),
('all-that-apply', 'What are your research field(s)?', 'f4a92aeee9f1d000'),
('long-answer', 'What do you think?', 'f4a92aeee9f1d000');

INSERT INTO `answer`(`question_id`, `value`) VALUES
(3, '<18'),
(3, '18-35'),
(3, 'old'),
(4, 'Physics'),
(4, 'An inferior science');
