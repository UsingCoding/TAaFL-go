import { generate } from "./escodegen/escodegen.js";
import { readLines } from "https://deno.land/std@0.76.0/io/bufio.ts";

async function fetchFromStdin() {
    for await (const line of readLines(Deno.stdin)) {
        return line;
    }
}

const code = await fetchFromStdin()
console.log(generate(code))