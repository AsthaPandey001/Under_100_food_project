export interface Category {
  id: string;
  name: string;
  startingPrice: number;
  image: string;
}

export const CATEGORIES: Category[] = [
  {
    id: "south-indian",
    name: "South Indian",
    startingPrice: 39,
    image: "/hero_food.jpg"
  },
  {
    id: "dosa-idli",
    name: "Dosa & Idli",
    startingPrice: 49,
    image: "/crispy_dosa.jpg"
  },
  {
    id: "rice-biryani",
    name: "Rice & Biryani",
    startingPrice: 79,
    image: "/veg_pulao.jpg"
  },
  {
    id: "rolls-wraps",
    name: "Rolls & Wraps",
    startingPrice: 69,
    image: "/kathi_roll.jpg"
  },
  {
    id: "snacks-chaat",
    name: "Snacks & Chaat",
    startingPrice: 39,
    image: "/samosa_plate.jpg"
  },
  {
    id: "north-indian",
    name: "North Indian Meals",
    startingPrice: 89,
    image: "/aloo_paratha.jpg"
  },
  {
    id: "momos",
    name: "Momos",
    startingPrice: 69,
    image: "/steamed_momos.jpg"
  },
  {
    id: "sandwiches",
    name: "Sandwiches",
    startingPrice: 59,
    image: "/bombay_sandwich.jpg"
  },
  {
    id: "parathas",
    name: "Parathas",
    startingPrice: 49,
    image: "/aloo_paratha.jpg"
  },
  {
    id: "beverages",
    name: "Beverages",
    startingPrice: 29,
    image: "/filter_coffee.jpg"
  }
];
