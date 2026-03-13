-- Delete duplicate scores, keeping only the best score per user
DELETE FROM "Score"
WHERE id NOT IN (
  SELECT DISTINCT ON ("userId") id
  FROM "Score"
  ORDER BY "userId", score DESC
);

-- CreateIndex
CREATE UNIQUE INDEX "Score_userId_key" ON "Score"("userId");
