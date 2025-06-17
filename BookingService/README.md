## Steps to setup the starter template

1. Clone the project

```
git clone https://github.com/singhsanket143/Express-Typescript-Starter-Project.git <ProjectName>
```

2. Move in to the folder structure

```
cd <ProjectName>
```

3. Install npm dependencies

```
npm i
```

4. Create a new .env file in the root directory
   and add the `PORT, DATABASE_URL, REDIS_SERVER_URL, lOCK_TTL` env variables

```
echo PORT=3000 >> .env
```

5. Migrate DB Changes (make sure valid **DATABASE_URL** provided previously in `.env` file

```
npx prisma migrate deploy
```

6. Generate prisma client / type

```
npx prisma generate
```

    7. 	make sure db server and redis server running,
    if so, then run below cmd get started 😇

```
npm run dev
```
