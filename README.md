![Size](https://img.shields.io/github/repo-size/totoledao/auction-house)
![Platform](https://img.shields.io/badge/platform-server-7F00FF)

[![GO][go-shield]][go-url]

<!-- PROJECT LOGO -->
<!-- <br /> -->
<p align="center">
  <a href="https://github.com/totoledao/auction-house">
    <!-- <img src="web\src\assets\logo.svg" alt="SpaceTime Logo" width="250"> -->
  </a>

  <p align="center">
    Real time auctioning!
  </p>
</p>

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary><h2 style="display: inline-block">Table of Contents</h2></summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
        <li><a href="#technologies">Technologies</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
        <!-- <li><a href="#usage">Usage</a></li> -->
      </ul>
    </li>    
    <!-- <li><a href="#license">License</a></li> -->
    <li><a href="#contact">Contact</a></li>    
  </ol>
</details>

[API documentation](APIDocumentation.md)

<!-- ABOUT THE PROJECT -->

## About The Project

Minimal implementation of a real-time auctioning house back-end service that enables users to list products, create auctions, place bids, and participate in live bidding events with real-time updates.

The main goal in developing this was to improve my proficiency in creating back-end services that implements real-time communication patterns, concurrent programming and authentication security features like CSRF protection and session management.

### Built With

- [Go][go-url]
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)

### Technologies

- Routing - [chi](https://github.com/go-chi/chi)<br>
- PostgreSQL Driver - [pgx](https://github.com/jackc/pgx/)<br>
- Session Management - [scs, pgxstore](https://github.com/alexedwards/scs)<br>
- Migrations - [tern](https://github.com/jackc/tern)<br>
- CSRF protection - [gorilla/csrf](https://github.com/gorilla/csrf)
- Websockets - [gorilla/websocket](https://github.com/gorilla/websocket)
- Accurate float handling - [shopspring/decimal](https://github.com/shopspring/decimal)

<!-- ### Features -->

<!-- ### Technical Goals -->

<!-- GETTING STARTED -->

## Getting Started

To get a local copy up and running follow these steps.

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/totoledao/auction-house.git
   ```
1. Create a .env file
1. Start the services
   ```sh
   docker compose up -d
   ```
1. Install and update dependencies
   ```sh
   go get -u ./...
   ```
1. Run migrations
   ```sh
   go install github.com/jackc/tern/v2@latest
   go run ./cmd/terndotenv
   ```
1. Start the server
   ```sh
   go run ./cmd/server
   ```

<!-- ## Usage -->

<!-- LICENSE -->

<!-- ## License

Distributed under the MIT License. See [`LICENSE`][license-url] for more information. -->

<!-- CONTACT -->

## Contact

Guilherme Toledo - guilherme-toledo@live.com

[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/guilhermemtoledo/)[![Instagram](https://img.shields.io/badge/Instagram-E4405F?style=for-the-badge&logo=instagram&logoColor=white)](https://www.instagram.com/totoledao)[![GitHub](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=whit)](https://www.github.com/totoledao)

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=0e76a8
[linkedin-url]: http://www.linkedin.com/in/guilhermemtoledo
[go-shield]: https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white
[go-url]: https://go.dev/
