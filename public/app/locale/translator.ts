import { vocabs } from "./fa-IR";

interface Vocabs {
  [key: string]: string;
}

export function Translator(vocabulary: Vocabs) {
  const dictionary = Object.entries(vocabulary).reduce(
    (dic, [key, value]) => ({
      ...dic,
      [key.toLowerCase()]: value,
    }),
    {} as Vocabs
  );

  return function translate(key = "") {
    if (typeof key !== "string") {
      console.log("delete -> " + key);
      console.log(new Error().stack);
      return key;
    }
    if (!key) {
      return key;
    }
    const term = key.toLowerCase();
    if (!dictionary[term]) {
      console.log("add -> " + key + " <-");
      console.log(new Error().stack);
      return key;
    }
    return dictionary[term] || key;
  };
}

export const translate = Translator(vocabs);
