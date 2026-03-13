-- Rename Score table to ScoreTetris
ALTER TABLE "Score" RENAME TO "ScoreTetris";

-- Rename constraints/indexes to match new table name
ALTER TABLE "ScoreTetris" RENAME CONSTRAINT "Score_pkey" TO "ScoreTetris_pkey";
ALTER TABLE "ScoreTetris" RENAME CONSTRAINT "Score_userId_fkey" TO "ScoreTetris_userId_fkey";
ALTER INDEX "Score_userId_key" RENAME TO "ScoreTetris_userId_key";

-- CreateTable ScoreWolfpack
CREATE TABLE "ScoreWolfpack" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "score" INTEGER NOT NULL,
    "waves" INTEGER NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "ScoreWolfpack_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "ScoreWolfpack_userId_key" ON "ScoreWolfpack"("userId");

-- AddForeignKey
ALTER TABLE "ScoreWolfpack" ADD CONSTRAINT "ScoreWolfpack_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;
