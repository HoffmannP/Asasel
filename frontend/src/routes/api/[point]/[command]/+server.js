import { json } from "@sveltejs/kit";

export async function POST({ params }) {
  console.log(params);
  return json([params.point, params.command]);
}
