export function titleCase(s: string) {
  return s
    .split(" ")
    .filter((c) => c.length > 0)
    .map((c) => c[0].toUpperCase() + c.substr(1).toLowerCase())
    .join(" ");
}
