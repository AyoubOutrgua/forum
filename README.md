# 🗨️ Forum — Zone01 Oujda Project

A complete web forum built using **Go**, **SQLite**, **HTML/CSS**, and **Docker**.  
This project focuses on backend and frontend fundamentals without relying on any JavaScript frameworks.  
It includes authentication, sessions, posting, commenting, likes/dislikes, categories, filters, and containerization.

---

## 📚 Table of Contents
- [Project Overview](#project-overview)
- [Features](#features)
- [Database (SQLite)](#database-sqlite)
- [Authentication](#authentication)
- [Communication Between Users](#communication-between-users)
- [Likes & Dislikes](#likes--dislikes)
- [Filtering System](#filtering-system)
- [Project Structure](#project-structure)
- [Docker](#docker)
- [Installation & Usage](#installation--usage)
- [Allowed Packages](#allowed-packages)
- [Learning Outcomes](#learning-outcomes)
- [Contributors](#contributors)

---

## 🧾 Project Overview

Your task is to build a functional **web forum** that supports:

- ✔️ User communication  
- ✔️ Posts with categories  
- ✔️ Comments  
- ✔️ Likes & dislikes  
- ✔️ Filtering posts  
- ✔️ Secure authentication system  
- ✔️ Running everything inside Docker  

The project must be written in **Go** and must use **SQLite** as the database engine.  
No frontend frameworks (React, Vue, Angular…) are allowed.

---

## ✨ Features

### 🔐 Authentication
- Registration using:
  - Email (must be unique)
  - Username
  - Password
- Login using cookies
- Only one active session per user
- Each session has an expiration date
- Password encryption (**bonus**) using bcrypt
- UUID sessions (**bonus**)
- Error handling:
  - Email already registered
  - Wrong password
  - Invalid credentials

---

### 📝 Communication Between Users
- Only registered users can:
  - Create posts  
  - Create comments  
- Posts can have **multiple categories**
- All users (even not logged in) can view:
  - Posts
  - Comments

---

### 👍 Likes & Dislikes
- Only logged-in users can like/dislike posts or comments
- Everyone can see:
  - Number of likes
  - Number of dislikes

---

### 🔎 Filtering System
Users can filter posts by:
- **Categories**
- **Posts created by the logged-in user**
- **Posts liked by the logged-in user**

The last two filters are available **only for authenticated users**.

---

## 🗄️ Database (SQLite)

You must use **SQLite** to store all forum data:

- Users  
- Sessions  
- Posts  
- Comments  
- Categories  
- Relations (post-category, likes, dislikes)

You must include at least:

- ✔️ One `CREATE` query  
- ✔️ One `INSERT` query  
- ✔️ One `SELECT` query  

It is recommended to design an **Entity Relationship Diagram (ERD)** before implementation.

---

## 📁 Project Structure (Suggested)

