# Campus Ambassador Platform - Backend API

A comprehensive Go-based backend API for a campus ambassador platform with gamification, campaigns, tasks, rewards, and engagement features.

## Table of Contents

- [Project Structure](#project-structure)
- [Technology Stack](#technology-stack)
- [Directory Structure](#directory-structure)
- [Core Files and Functions](#core-files-and-functions)
- [API Routes](#api-routes)
- [Database Models](#database-models)
- [Setup and Installation](#setup-and-installation)

## Project Structure

```
backend/
├── cmd/server/          # Main application code and HTTP handlers
├── internal/            # Internal packages (db, store, services, env)
│   ├── db/             # Database connection and seeding
│   ├── env/            # Environment variable management
│   ├── services/       # Business logic services
│   └── store/          # Database models and store functions
├── migrations/         # Database migration files
├── ws/                # WebSocket hub implementation
├── graph/             # GraphQL schema and resolvers (if exists)
├── bin/               # Compiled binaries
├── docker-compose.yml # Docker configuration
├── Dockerfile         # Docker image definition
├── go.mod            # Go module dependencies
├── go.sum            # Go module checksums
├── gqlgen.yml        # GraphQL code generation config
└── Makefile          # Build and run commands
```

## Technology Stack

- **Language**: Go 1.24.3
- **Web Framework**: Chi Router
- **Database**: PostgreSQL with GORM
- **GraphQL**: gqlgen
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **WebSocket**: Gorilla WebSocket
- **Validation**: go-playground/validator/v10

## Directory Structure

### `/cmd/server/` - Application Entry Point and Handlers

Contains all HTTP handlers, middleware, and the main application entry point.

#### Core Files

##### `main.go`
**Purpose**: Application entry point and server initialization

**Functions**:
- `main()` - Initializes the server, sets up routes, middleware, GraphQL, REST API, and WebSocket handlers

##### `rest.go`
**Purpose**: REST API route definitions

**Functions**:
- `setupREST(r chi.Router, db *gorm.DB)` - Configures all REST API routes
- `healthHandler(w http.ResponseWriter, r *http.Request)` - Health check endpoint

**Routes Configured**:
- `/api/v1/health` - Health check
- `/api/v1/auth/*` - Authentication routes
- `/api/v1/users/*` - User management routes
- `/api/v1/tasks/*` - Task management routes
- `/api/v1/submissions/*` - Submission routes
- `/api/v1/campaigns/*` - Campaign routes
- `/api/v1/xp/*` - XP and gamification routes
- `/api/v1/rewards/*` - Rewards routes
- `/api/v1/referrals/*` - Referral routes
- `/api/v1/engagement/*` - Engagement features (flash challenges, trivia, battles)
- `/api/v1/wars/*` - Campus wars routes
- `/api/v1/surveys/*` - Survey routes
- `/api/v1/notifications/*` - Notification routes

##### `json-api.go`
**Purpose**: JSON request/response utilities

**Functions**:
- `writeJSON(w http.ResponseWriter, status int, data any) error` - Writes JSON response
- `readJSON(w http.ResponseWriter, r *http.Request, data any) error` - Reads and parses JSON request
- `writeJSONError(w http.ResponseWriter, status int, message string)` - Writes JSON error response
- `jsonResponse(w http.ResponseWriter, status int, data any) error` - Writes JSON envelope response

**Variables**:
- `Validate *validator.Validate` - Request validation instance

##### `error-api.go`
**Purpose**: HTTP error response handlers

**Functions**:
- `internalServerError(w http.ResponseWriter, r *http.Request, err error)` - Returns 500 error
- `badRequestResponse(w http.ResponseWriter, r *http.Request, err error)` - Returns 400 error
- `notFoundResponse(w http.ResponseWriter, r *http.Request, err error)` - Returns 404 error
- `conflictResponse(w http.ResponseWriter, r *http.Request, err error)` - Returns 409 error
- `unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error)` - Returns 401 error

##### `middleware.go`
**Purpose**: HTTP middleware functions

**Functions**:
- `RequireAuth(db *gorm.DB) func(http.Handler) http.Handler` - Authentication middleware
- `GetUserFromContext(r *http.Request) (*store.User, bool)` - Extracts user from request context

**Context Keys**:
- `userContextKey` - Context key for storing authenticated user

##### `auth.go`
**Purpose**: Authentication and user registration

**Functions**:
- `generateToken(user *store.User) (string, string, error)` - Generates JWT access and refresh tokens
- `loginHandler(db *gorm.DB) http.HandlerFunc` - User login handler
- `registerHandler(db *gorm.DB) http.HandlerFunc` - User registration handler
- `refreshTokenHandler(db *gorm.DB) http.HandlerFunc` - Token refresh handler
- `logoutHandler(db *gorm.DB) http.HandlerFunc` - User logout handler
- `forgotPasswordHandler(db *gorm.DB) http.HandlerFunc` - Password reset request handler
- `resetPasswordHandler(db *gorm.DB) http.HandlerFunc` - Password reset handler
- `verifyEmailHandler(db *gorm.DB) http.HandlerFunc` - Email verification handler
- `generateReferralCode() string` - Generates random referral code
- `generateSecureToken(length int) string` - Generates secure random token
- `extractToken(authHeader string) string` - Extracts Bearer token from header
- `updateUserStreak(db *gorm.DB, userID uint, streakType string) error` - Updates user login streak

**Request Types**:
- `LoginRequest` - Login request payload
- `RegisterRequest` - Registration request payload
- `RefreshTokenRequest` - Token refresh request
- `ForgotPasswordRequest` - Password reset request
- `ResetPasswordRequest` - Password reset with token

##### `user.go`
**Purpose**: User profile management

**Functions**:
- `getCurrentUserProfileHandler(db *gorm.DB) http.HandlerFunc` - Get current user profile
- `updateUserProfileHandler(db *gorm.DB) http.HandlerFunc` - Update user profile
- `updateAvatarHandler(db *gorm.DB) http.HandlerFunc` - Update user avatar
- `uploadResumeHandler(db *gorm.DB) http.HandlerFunc` - Upload user resume
- `getUserCertificatesHandler(db *gorm.DB) http.HandlerFunc` - Get user certificates
- `downloadCertificateHandler(db *gorm.DB) http.HandlerFunc` - Download certificate

**Request Types**:
- `UpdateProfileRequest` - Profile update payload

##### `task.go`
**Purpose**: Task management

**Functions**:
- `getTasksHandler(db *gorm.DB) http.HandlerFunc` - Get tasks with filters
- `getTaskHandler(db *gorm.DB) http.HandlerFunc` - Get single task
- `getAssignedTasksHandler(db *gorm.DB) http.HandlerFunc` - Get assigned tasks
- `getAvailableTasksHandler(db *gorm.DB) http.HandlerFunc` - Get available tasks

##### `submission.go`
**Purpose**: Task submission management

**Functions**:
- `getSubmissionsHandler(db *gorm.DB) http.HandlerFunc` - Get user submissions
- `createSubmissionHandler(db *gorm.DB) http.HandlerFunc` - Create new submission
- `getSubmissionHandler(db *gorm.DB) http.HandlerFunc` - Get single submission
- `updateSubmissionHandler(db *gorm.DB) http.HandlerFunc` - Update submission
- `deleteSubmissionHandler(db *gorm.DB) http.HandlerFunc` - Delete submission
- `getSubmissionProofHandler(db *gorm.DB) http.HandlerFunc` - Get submission proof

**Response Types**:
- `SubmissionResponse` - Submission data structure

##### `campaign.go`
**Purpose**: Campaign management

**Functions**:
- `getCampaignsHandler(db *gorm.DB) http.HandlerFunc` - Get campaigns
- `getCampaignHandler(db *gorm.DB) http.HandlerFunc` - Get single campaign
- `joinCampaignHandler(db *gorm.DB) http.HandlerFunc` - Join campaign
- `getCampaignTasksHandler(db *gorm.DB) http.HandlerFunc` - Get campaign tasks
- `getCampaignLeaderboardHandler(db *gorm.DB) http.HandlerFunc` - Get campaign leaderboard

##### `gamification.go`
**Purpose**: Gamification features (XP, Levels, Badges, Streaks, Spin Wheel)

**Functions**:
- `getXPTransactionsHandler(db *gorm.DB) http.HandlerFunc` - Get XP transaction history
- `awardXPHandler(db *gorm.DB) http.HandlerFunc` - Award XP (admin only)
- `getLevelsHandler(db *gorm.DB) http.HandlerFunc` - Get all levels
- `getCurrentLevelHandler(db *gorm.DB) http.HandlerFunc` - Get current user level
- `getBadgesHandler(db *gorm.DB) http.HandlerFunc` - Get all badges
- `getBadgeHandler(db *gorm.DB) http.HandlerFunc` - Get single badge
- `getUserBadgesHandler(db *gorm.DB) http.HandlerFunc` - Get user badges
- `getStreakHandler(db *gorm.DB) http.HandlerFunc` - Get user streak
- `logStreakHandler(db *gorm.DB) http.HandlerFunc` - Log streak activity
- `getSpinWheelHandler(db *gorm.DB) http.HandlerFunc` - Get spin wheel config
- `spinWheelHandler(db *gorm.DB) http.HandlerFunc` - Spin the wheel
- `getSpinHistoryHandler(db *gorm.DB) http.HandlerFunc` - Get spin history

##### `engagement.go`
**Purpose**: Engagement features (Flash Challenges, Trivia, Mystery Boxes, Secret Codes, Content Battles)

**Functions**:
- `getActiveFlashChallengesHandler(db *gorm.DB) http.HandlerFunc` - Get active flash challenges
- `participateFlashChallengeHandler(db *gorm.DB) http.HandlerFunc` - Participate in flash challenge
- `getActiveTriviaHandler(db *gorm.DB) http.HandlerFunc` - Get active trivia tournaments
- `startTriviaHandler(db *gorm.DB) http.HandlerFunc` - Start trivia session
- `submitTriviaAnswersHandler(db *gorm.DB) http.HandlerFunc` - Submit trivia answers
- `getMysteryBoxesHandler(db *gorm.DB) http.HandlerFunc` - Get available mystery boxes
- `openMysteryBoxHandler(db *gorm.DB) http.HandlerFunc` - Open mystery box
- `redeemSecretCodeHandler(db *gorm.DB) http.HandlerFunc` - Redeem secret code
- `getWeeklyVibeChallengeHandler(db *gorm.DB) http.HandlerFunc` - Get weekly challenge
- `submitWeeklyVibeHandler(db *gorm.DB) http.HandlerFunc` - Submit weekly challenge entry
- `getWeeklyChallengeSubmissionsHandler(db *gorm.DB) http.HandlerFunc` - Get weekly submissions
- `voteWeeklyChallengeHandler(db *gorm.DB) http.HandlerFunc` - Vote on weekly submission
- `getActiveBattlesHandler(db *gorm.DB) http.HandlerFunc` - Get active content battles
- `submitBattleHandler(db *gorm.DB) http.HandlerFunc` - Submit battle entry
- `voteBattleHandler(db *gorm.DB) http.HandlerFunc` - Vote on battle submission

##### `rewards.go`
**Purpose**: Rewards and redemptions

**Functions**:
- `getRewardsHandler(db *gorm.DB) http.HandlerFunc` - Get available rewards
- `getRewardHandler(db *gorm.DB) http.HandlerFunc` - Get single reward
- `redeemRewardHandler(db *gorm.DB) http.HandlerFunc` - Redeem reward
- `getRewardRedemptionsHandler(db *gorm.DB) http.HandlerFunc` - Get user redemptions
- `updateRedemptionStatusHandler(db *gorm.DB) http.HandlerFunc` - Update redemption status (admin)
- `cancelRedemptionHandler(db *gorm.DB) http.HandlerFunc` - Cancel redemption

##### `referral.go`
**Purpose**: Referral system

**Functions**:
- `getReferralsHandler(db *gorm.DB) http.HandlerFunc` - Get user referrals
- `getReferralCodeHandler(db *gorm.DB) http.HandlerFunc` - Get user referral code
- `getReferralInvitesHandler(db *gorm.DB) http.HandlerFunc` - Get referral invites
- `sendReferralInviteHandler(db *gorm.DB) http.HandlerFunc` - Send referral invite

##### `wars.go`
**Purpose**: Campus Wars feature

**Functions**:
- `getActiveWarsHandler(db *gorm.DB) http.HandlerFunc` - Get active campus wars
- `getWarHandler(db *gorm.DB) http.HandlerFunc` - Get single war details
- `getWarParticipantsHandler(db *gorm.DB) http.HandlerFunc` - Get war participants
- `getWarLeaderboardHandler(db *gorm.DB) http.HandlerFunc` - Get war leaderboard

##### `survey.go`
**Purpose**: Survey management

**Functions**:
- `getAvailableSurveysHandler(db *gorm.DB) http.HandlerFunc` - Get available surveys
- `getSurveyHandler(db *gorm.DB) http.HandlerFunc` - Get survey details
- `submitSurveyHandler(db *gorm.DB) http.HandlerFunc` - Submit survey response
- `getSurveyResponsesHandler(db *gorm.DB) http.HandlerFunc` - Get user survey responses

##### `notification.go`
**Purpose**: Notification management

**Functions**:
- `getNotificationsHandler(db *gorm.DB) http.HandlerFunc` - Get user notifications
- `getUnreadNotificationsCountHandler(db *gorm.DB) http.HandlerFunc` - Get unread count
- `markNotificationReadHandler(db *gorm.DB) http.HandlerFunc` - Mark notification as read
- `markAllNotificationsReadHandler(db *gorm.DB) http.HandlerFunc` - Mark all as read
- `deleteNotificationHandler(db *gorm.DB) http.HandlerFunc` - Delete notification

##### `websocket.go`
**Purpose**: WebSocket connection handling

**Functions**:
- `serveWS(hub *ws.Hub, w http.ResponseWriter, r *http.Request)` - Upgrade HTTP to WebSocket
- `clientReader(hub *ws.Hub, client *ws.Client)` - Read messages from client
- `clientWriter(hub *ws.Hub, client *ws.Client)` - Write messages to client

**Variables**:
- `upgrader websocket.Upgrader` - WebSocket connection upgrader

##### `graphql.go`
**Purpose**: GraphQL server setup

**Functions**:
- `graphqlHandler(db *gorm.DB, hub *ws.Hub) http.Handler` - Creates GraphQL handler
- `graphqlPlayground() http.Handler` - GraphQL Playground handler
- `mountGraphQL(mux *http.ServeMux, db *gorm.DB, hub *ws.Hub)` - Mounts GraphQL endpoints

### `/internal/` - Internal Packages

#### `/internal/db/` - Database Management

##### `db.go`
**Purpose**: Database connection management

**Functions**:
- `Connect() *gorm.DB` - Establishes PostgreSQL connection using environment variables

**Variables**:
- `DB *gorm.DB` - Global database instance

**Environment Variables Used**:
- `DB_HOST` - Database host
- `DB_USER` - Database user
- `DB_PASS` - Database password
- `DB_NAME` - Database name
- `DB_PORT` - Database port

##### `seed.go`
**Purpose**: Database seeding (if implemented)

#### `/internal/env/` - Environment Configuration

##### `env.go`
**Purpose**: Environment variable management

**Functions**:
- `Load()` - Loads environment variables from .env file
- `Get(key, fallback string) string` - Gets environment variable with fallback

#### `/internal/services/` - Business Logic Services

##### `notifier.go`
**Purpose**: Notification service (if implemented)

#### `/internal/store/` - Database Models and Store Functions

Contains GORM models and database access functions for all entities.

##### `user.go`
**Models**: `User`, `UserSession`

**Functions**:
- `CreateUser(db *gorm.DB, u *User) error`
- `GetUserByID(db *gorm.DB, id uint) (*User, error)`
- `GetUserByUUID(db *gorm.DB, uuid string) (*User, error)`
- `GetUserByEmail(db *gorm.DB, email string) (*User, error)`
- `UpdateUser(db *gorm.DB, u *User) error`
- `DeleteUser(db *gorm.DB, id uint) error`
- `CreateSession(db *gorm.DB, s *UserSession) error`
- `GetSessionByToken(db *gorm.DB, token string) (*UserSession, error)`
- `DeleteSessionByToken(db *gorm.DB, token string) error`
- `DeleteExpiredSessions(db *gorm.DB) error`
- `GetUserByReferralCode(db *gorm.DB, referralCode string) (*User, error)`
- `GetUserStreak(db *gorm.DB, userID uint, streakType string) (*UserStreak, error)`
- `UpdateUserStreak(db *gorm.DB, streak *UserStreak) error`
- `GetUserWithRelations(db *gorm.DB, userID uint) (*User, error)`
- `GetUserBadgeCount(db *gorm.DB, userID uint) (int, error)`
- `GetUserCertificates(db *gorm.DB, userID uint) ([]Certificate, error)`

##### `campaign_task.go`
**Models**: `Campaign`, `Task`, `TaskAssignment`

**Functions**:
- `CreateCampaign(db *gorm.DB, campaign *Campaign) error`
- `GetCampaignByID(db *gorm.DB, id uint) (*Campaign, error)`
- `CreateTask(db *gorm.DB, task *Task) error`
- `GetTaskByID(db *gorm.DB, id uint) (*Task, error)`
- `CreateTaskAssignment(db *gorm.DB, assignment *TaskAssignment) error`
- `GetTaskAssignmentByID(db *gorm.DB, id uint) (*TaskAssignment, error)`
- `UpdateCampaign(db *gorm.DB, campaign *Campaign) error`
- `GetSubmissionsByUserAndTask(db *gorm.DB, userID uint, taskID uint) ([]Submission, error)`

##### `submission.go`
**Models**: `Submission`, `SubmissionMedia`

**Functions**:
- `CreateSubmission(db *gorm.DB, submission *Submission) error`
- `GetSubmissionByID(db *gorm.DB, id uint) (*Submission, error)`
- `CreateSubmissionMedia(db *gorm.DB, media *SubmissionMedia) error`
- `GetSubmissionMediaByID(db *gorm.DB, id uint) (*SubmissionMedia, error)`
- `UpdateSubmission(db *gorm.DB, submission *Submission) error`
- `DeleteSubmission(db *gorm.DB, id uint) error`
- `GetSubmissionsByUserAndTasks(db *gorm.DB, userID uint, taskIDs []uint) ([]Submission, error)`
- `GetTaskAssignmentsByUser(db *gorm.DB, userID uint) ([]TaskAssignment, error)`

##### `levels_badges.go`
**Models**: `Level`, `Badge`, `ProfileSkin`

**Functions**:
- `CreateLevel(db *gorm.DB, level *Level) error`
- `GetLevelByID(db *gorm.DB, id uint) (*Level, error)`
- `CreateBadge(db *gorm.DB, badge *Badge) error`
- `GetBadgeByID(db *gorm.DB, id uint) (*Badge, error)`
- `CreateProfileSkin(db *gorm.DB, skin *ProfileSkin) error`
- `GetProfileSkinByID(db *gorm.DB, id uint) (*ProfileSkin, error)`

##### `user_badges.go`
**Models**: `UserBadge`

**Functions**:
- `CreateUserBadge(db *gorm.DB, userBadge *UserBadge) error`
- `GetUserBadgeByID(db *gorm.DB, id uint) (*UserBadge, error)`
- `GetUserBadges(db *gorm.DB, userID uint) ([]UserBadge, error)`

##### `gamification.go`
**Models**: `UserStreak`, `StreakLog`

**Functions**:
- `CreateUserStreak(db *gorm.DB, streak *UserStreak) error`
- `GetUserStreakByID(db *gorm.DB, id uint) (*UserStreak, error)`
- `CreateStreakLog(db *gorm.DB, log *StreakLog) error`
- `GetStreakLogByID(db *gorm.DB, id uint) (*StreakLog, error)`

##### `xp_transaction.go`
**Models**: `XPTransaction`

**Functions**:
- `CreateXPTransaction(db *gorm.DB, transaction *XPTransaction) error`
- `GetXPTransactionByID(db *gorm.DB, id uint) (*XPTransaction, error)`

##### `reward_store.go`
**Models**: `RewardStore`, `UserReward`

**Functions**:
- `CreateRewardStore(db *gorm.DB, reward *RewardStore) error`
- `GetRewardStoreByID(db *gorm.DB, id uint) (*RewardStore, error)`
- `CreateUserReward(db *gorm.DB, userReward *UserReward) error`
- `GetUserRewardByID(db *gorm.DB, id uint) (*UserReward, error)`

##### `referrals.go`
**Models**: `Referral`

**Functions**:
- `CreateReferral(db *gorm.DB, referral *Referral) error`

##### `college_state.go`
**Models**: `State`, `College`

**Functions**:
- `CreateState(db *gorm.DB, state *State) error`
- `GetStateByID(db *gorm.DB, id uint) (*State, error)`
- `CreateCollege(db *gorm.DB, college *College) error`
- `GetCollegeByID(db *gorm.DB, id uint) (*College, error)`

##### `notifications.go`
**Models**: `Notification`

**Functions**: (Notification store functions)

##### `flash_challenge.go`
**Models**: `FlashChallenge`

**Functions**: (Flash challenge store functions)

##### `trivia_tournament.go`
**Models**: `TriviaTournament`, `TriviaParticipant`

**Functions**: (Trivia store functions)

##### `campus_wars.go`
**Models**: `CampusWar`, `WarParticipant`

**Functions**: (Campus wars store functions)

##### `content_battles.go`
**Models**: `ContentBattle`, `BattleSubmission`

**Functions**: (Content battles store functions)

##### `spin_wheel.go`
**Models**: `SpinWheel`, `SpinWheelItem`, `UserSpin`

**Functions**: (Spin wheel store functions)

##### `secret_code.go`
**Models**: `SecretCode`, `UserSecretCode`

**Functions**: (Secret code store functions)

##### `badge_bingo.go`
**Models**: `BadgeBingo`, `UserBingoProgress`

**Functions**:
- `CreateBadgeBingo(db *gorm.DB, bingo *BadgeBingo) error`
- `GetBadgeBingoByID(db *gorm.DB, id uint) (*BadgeBingo, error)`
- `CreateUserBingoProgress(db *gorm.DB, progress *UserBingoProgress) error`
- `GetUserBingoProgressByID(db *gorm.DB, id uint) (*UserBingoProgress, error)`

##### `survey.go`
**Models**: `Survey`, `SurveyResponse`

**Functions**:
- `CreateSurvey(db *gorm.DB, survey *Survey) error`
- `GetSurveyByID(db *gorm.DB, id uint) (*Survey, error)`
- `CreateSurveyResponse(db *gorm.DB, response *SurveyResponse) error`
- `GetSurveyResponseByID(db *gorm.DB, id uint) (*SurveyResponse, error)`

##### `certificates.go`
**Models**: `Certificate`

**Functions**: (Certificate store functions)

##### `quest.go`
**Models**: `Quest`, `QuestProgress`

**Functions**: (Quest store functions)

##### `leaderboard.go`
**Models**: `Leaderboard`, `LeaderboardEntry`

**Functions**: (Leaderboard store functions)

##### `activity_logs.go`
**Models**: `ActivityLog`

**Functions**: (Activity log store functions)

##### `admin_system.go`
**Models**: `AdminAction`, `AdminUser`

**Functions**: (Admin system store functions)

### `/ws/` - WebSocket Hub

##### `hub.go`
**Purpose**: WebSocket hub for real-time communication

**Types**:
- `Client` - WebSocket client connection
- `Hub` - WebSocket hub managing all connections

**Functions**:
- `NewHub() *Hub` - Creates new WebSocket hub
- `(h *Hub) Run()` - Runs the hub's main loop

**Channels**:
- `Clients map[*Client]bool` - Active client connections
- `Broadcast chan []byte` - Broadcast message channel
- `Register chan *Client` - Client registration channel
- `Unregister chan *Client` - Client unregistration channel

### `/migrations/` - Database Migrations

Contains SQL migration files for database schema:
- User tables (000001)
- College/State tables (000002)
- Levels/Badges tables (000003)
- Campaign/Task tables (000004)
- User badges (000005)
- Submissions (000006)
- Gamification (000007)
- Spin wheel (000008)
- Trivia tournament (000009)
- Campus wars (000010)
- Secret codes (000011)
- Flash challenges (000012)
- Badge bingo (000013)
- Content battles (000014)
- Quest table (000015)
- Reward store (000016)
- Certificates (000017)
- Leaderboard (000018)
- XP transactions (000019)
- Activity logs (000020)
- Referrals (000021)
- Notifications (000022)
- Surveys (000023)
- Admin system (000024)
- Indexing (000025)
- Views (000026)
- Triggers (000027)

## API Routes

### Public Routes
- `GET /api/v1/health` - Health check

### Authentication Routes (`/api/v1/auth`)
- `POST /login` - User login
- `POST /register` - User registration
- `POST /refresh` - Refresh access token
- `POST /logout` - User logout
- `POST /forgot-password` - Request password reset
- `POST /reset-password` - Reset password with token
- `GET /verify-email/{token}` - Verify email address

### Protected Routes (Require Authentication)

#### User Routes (`/api/v1/users`)
- `GET /me` - Get current user profile
- `PUT /me` - Update user profile
- `PATCH /me/avatar` - Update avatar
- `POST /me/resume` - Upload resume
- `GET /me/certificates` - Get user certificates
- `GET /me/certificates/{id}/download` - Download certificate

#### Task Routes (`/api/v1/tasks`)
- `GET /` - Get tasks (with filters)
- `GET /{id}` - Get single task
- `GET /assigned` - Get assigned tasks
- `GET /available` - Get available tasks

#### Submission Routes (`/api/v1/submissions`)
- `GET /` - Get user submissions
- `POST /` - Create submission
- `GET /{id}` - Get single submission
- `PUT /{id}` - Update submission
- `DELETE /{id}` - Delete submission
- `GET /{id}/proof` - Get submission proof

#### Campaign Routes (`/api/v1/campaigns`)
- `GET /` - Get campaigns
- `GET /{id}` - Get single campaign
- `POST /{id}/join` - Join campaign
- `GET /{id}/tasks` - Get campaign tasks
- `GET /{id}/leaderboard` - Get campaign leaderboard

#### Gamification Routes

**XP (`/api/v1/xp`)**
- `GET /transactions` - Get XP transactions
- `POST /award` - Award XP (admin)

**Levels (`/api/v1/levels`)**
- `GET /` - Get all levels
- `GET /current` - Get current user level

**Badges (`/api/v1/badges`)**
- `GET /` - Get all badges
- `GET /{id}` - Get single badge
- `GET /me` - Get user badges

**Streaks (`/api/v1/streaks`)**
- `GET /` - Get user streak
- `POST /log` - Log streak activity

**Spin Wheel (`/api/v1/spin-wheel`)**
- `GET /` - Get spin wheel config
- `POST /spin` - Spin the wheel
- `GET /history` - Get spin history

#### Engagement Routes

**Flash Challenges (`/api/v1/flash-challenges`)**
- `GET /active` - Get active challenges
- `POST /{id}/participate` - Participate in challenge

**Trivia (`/api/v1/trivia`)**
- `GET /active` - Get active trivia
- `POST /{id}/start` - Start trivia session
- `POST /{id}/submit-answers` - Submit answers

**Mystery Boxes (`/api/v1/mystery-boxes`)**
- `GET /` - Get available boxes
- `POST /{id}/open` - Open mystery box

**Secret Codes (`/api/v1/secret-codes`)**
- `POST /redeem/{code}` - Redeem secret code

**Weekly Challenge (`/api/v1/weekly-challenge`)**
- `GET /current` - Get current weekly challenge
- `POST /submit` - Submit weekly entry
- `GET /submissions` - Get submissions
- `POST /vote/{submissionId}` - Vote on submission

**Battles (`/api/v1/battles`)**
- `GET /active` - Get active battles
- `POST /{id}/submit` - Submit battle entry
- `POST /{id}/vote/{submissionId}` - Vote on battle entry

#### Rewards (`/api/v1/rewards`)
- `GET /` - Get available rewards
- `GET /{id}` - Get single reward
- `POST /{id}/redeem` - Redeem reward
- `GET /redemptions` - Get user redemptions

#### Referrals (`/api/v1/referrals`)
- `GET /` - Get user referrals
- `GET /code` - Get referral code
- `GET /invites` - Get referral invites
- `POST /invite` - Send referral invite

#### Campus Wars (`/api/v1/wars`)
- `GET /active` - Get active wars
- `GET /{id}` - Get war details
- `GET /{id}/participants` - Get participants
- `GET /{id}/leaderboard` - Get war leaderboard

#### Surveys (`/api/v1/surveys`)
- `GET /available` - Get available surveys
- `GET /{id}` - Get survey details
- `POST /{id}/submit` - Submit survey response
- `GET /responses` - Get user responses

#### Notifications (`/api/v1/notifications`)
- `GET /` - Get notifications
- `GET /unread-count` - Get unread count
- `PUT /{id}/read` - Mark as read
- `PUT /read-all` - Mark all as read
- `DELETE /{id}` - Delete notification

### GraphQL Endpoints
- `POST /graphql` - GraphQL endpoint
- `GET /playground` - GraphQL Playground (if enabled)

### WebSocket
- `WS /ws` - WebSocket connection endpoint

## Database Models

### Core Models

#### User
- ID, UUID, Email, PasswordHash
- FirstName, LastName, Phone
- Role, CollegeID, StateID
- ReferralCode, ReferredBy
- XP, LevelID, StreakCount
- TotalSubmissions, ApprovedSubmissions, WinRate
- AvatarURL, ResumeURL
- IsActive, CreatedAt, UpdatedAt

#### Campaign
- ID, UUID, Title, Description
- CampaignType, Category
- BannerImageURL
- StartDate, EndDate
- MaxParticipants, CurrentParticipants
- Status, Priority
- CreatedBy, IsLimitedEdition
- Metadata, CreatedAt, UpdatedAt

#### Task
- ID, UUID, CampaignID
- Title, Description
- TaskType, ProofType
- XPReward, CoinReward
- DurationHours, Priority
- IsActive, SubmissionInstructions
- CreatedBy, CreatedAt, UpdatedAt

#### Submission
- ID, UUID, TaskID, UserID, CampaignID
- ProofType, ProofURL, ProofText
- Status, SubmissionStage
- ReviewedBy, ReviewedAt, ReviewComments
- XPAwarded, CoinsAwarded, Score
- SubmittedAt, UpdatedAt

### Gamification Models

#### Level
- ID, Name, RankOrder
- MinXP, MaxXP
- BadgeURL, Description
- UnlockableFeatures

#### Badge
- ID, Name, Description
- BadgeType, Category
- ImageURL, XPReward
- CriteriaType, CriteriaValue
- IsSecret, IsLimitedEdition

#### UserBadge
- ID, UserID, BadgeID
- EarnedAt

#### XPTransaction
- ID, UserID
- TransactionType, Amount, BalanceAfter
- SourceType, SourceID
- Description, Metadata
- CreatedAt

#### UserStreak
- ID, UserID, StreakType
- CurrentStreak, LongestStreak
- LastActivityDate, TotalDays

### Reward Models

#### RewardStore
- ID, Name, Description
- RewardType, Category
- ImageURL
- XPCost, CoinCost, CashCost
- QuantityAvailable, QuantitySold
- IsFeatured, IsActive
- ValidityDays

#### UserReward
- ID, UserID, RewardID
- RedemptionCode, Status
- XPPaid, CoinsPaid, CashPaid
- ShippingAddress, TrackingNumber
- ClaimedAt, DeliveredAt

### Engagement Models

#### FlashChallenge
- ID, Title, Description
- ChallengeType, DurationHours
- XPReward, MaxParticipants
- StartTime, EndTime
- Status

#### TriviaTournament
- ID, Title, Description
- TournamentType, QuestionCount
- StartDate, EndDate
- DurationMinutes, MaxParticipants
- EntryFeeXp, Status

#### ContentBattle
- ID, Title, Description
- BattleType, Theme
- SubmissionDeadline
- VotingStart, VotingEnd
- MaxParticipants, Status

#### CampusWar
- ID, Name, Description
- WarType
- StartDate, EndDate
- Status

### Other Models

#### Referral
- ID, ReferrerID, ReferredEmail
- ReferredUserID, Status
- XPAwarded, XPAwardedToReferred
- ConversionStage

#### Notification
- ID, UserID
- NotificationType
- Title, Message, Data
- IsRead, IsActionable
- ActionURL
- SentAt, ReadAt

#### Survey
- ID, Title, Description
- SurveyType, Questions
- XPReward, IsActive
- StartDate, EndDate
- CreatedBy

## Setup and Installation

### Prerequisites
- Go 1.24.3 or higher
- PostgreSQL database
- Environment variables configured

### Environment Variables

Create a `.env` file with:
```
DB_HOST=localhost
DB_USER=postgres
DB_PASS=postgres
DB_NAME=yourapp
DB_PORT=5432
SERVER_PORT=8080
JWT_SECRET=your-secret-key
JWT_REFRESH=your-refresh-secret
```

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run migrations:
   ```bash
   make migrate-up
   ```

4. (Optional) Seed database:
   ```bash
   make seed
   ```

5. Run the server:
   ```bash
   make run
   ```
   Or build and run:
   ```bash
   make build
   ./bin/server
   ```

### Docker Setup

```bash
make docker
```

## Makefile Commands

- `make run` - Run the server
- `make build` - Build the binary
- `make docker` - Run with Docker Compose
- `make migrate-up` - Run database migrations
- `make seed` - Seed the database

## Architecture Notes

- **Layered Architecture**: Handlers → Store → Database
- **Middleware**: Authentication, CORS, Logging, Recovery
- **Dual API**: REST and GraphQL support
- **Real-time**: WebSocket support for live updates
- **Gamification**: Comprehensive XP, levels, badges, streaks system
- **Engagement**: Multiple engagement features (challenges, trivia, battles, wars)
- **Modular**: Clear separation of concerns with internal packages

