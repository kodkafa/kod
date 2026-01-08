args <- commandArgs(trailingOnly = TRUE)
name_val <- "RUser"

if (length(args) > 1 && args[1] == "--name") {
  name_val <- args[2]
}

cat("╔════════════════════════════════════╗\n")
cat(sprintf("║  Name: %-27s ║\n", name_val))
cat("║  Hello from KODKAFA R Plugin!      ║\n")
cat("╚════════════════════════════════════╝\n")
