Live Link  :- https://exploding-kittens-nine.vercel.app/


# 😸 Exploding Kitten Game

Welcome to the Exploding Kitten game! This is an online single-player card game built using React, Redux, Golang, and Redis. The objective of the game is to draw cards from the deck without drawing the Exploding bomb card. 

## 🎮 How to Play

1. **Start Game**: Enter your username and start the game.
2. **Draw Cards**: Click on the deck to draw a card. Keep drawing until you either win by drawing all non-exploding kitten cards or lose by drawing an exploding bomb.
3. **Win or Lose**: If you draw all non-exploding cards, you win! Your win will be recorded in the leaderboard. If you draw an exploding kitten, you lose.

## 🚀 Technologies Used

- **Frontend**: React, Redux
- **Backend**: Golang
- **Database**: Redis

### 🛠️ Setup
 
Steps to Run the Application

1. Clone the Repository
Clone the repository to your local system using the following command:

git clone https://github.com/pragya2003/assignment-exploding-kitten-game.git
cd https://github.com/pragya2003/assignment-exploding-kitten-game.git

2. Install Dependencies
Use the package manager to install all the required dependencies:

npm install
or
yarn install

3. Setup Environment Variables
If the project requires environment variables:
Locate the .env.example file in the project directory.
Copy it to create a .env file:

cp .env.example .env

Fill in the required values as specified.

4. Start the Development Server
Run the development server:

npm start
or

yarn start
Open the application in your browser at http://localhost:3000.

5. Build for Production
To create a production build, run:

npm run build
or

yarn build
This will create a build folder containing optimized files.

6. Serve the Production Build (Optional)
To preview the production build locally:

npm install -g serve
serve -s build
Access the application at http://localhost:5000.

7. Run Tests (If Applicable)
To run tests:

npm test
or

yarn test


## 📊 Leaderboard

The leaderboard records the number of games won by each player. One game won equals one point.


