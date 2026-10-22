package seeder

import (
	"context"
	"fmt"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/city"
	"github.com/mikestefanello/pagoda/ent/restaurant"
	"github.com/mikestefanello/pagoda/ent/user"
	"golang.org/x/crypto/bcrypt"
)

// SeedFoodDiscovery seeds Bangalore city, admin/normal users, and test restaurants
func SeedFoodDiscovery(ctx context.Context, client *ent.Client) error {
	// 1. Seed Users (Admin, Customer, and Merchant roles)
	hashedAdminPassword, _ := bcrypt.GenerateFromPassword([]byte("Admin@123456"), bcrypt.DefaultCost)
	adminUser, err := client.User.Query().Where(user.EmailEQ("admin@fooddiscovery.com")).Only(ctx)
	if ent.IsNotFound(err) {
		adminUser, err = client.User.Create().
			SetName("Food Discovery Admin").
			SetEmail("admin@fooddiscovery.com").
			SetPassword(string(hashedAdminPassword)).
			SetVerified(true).
			SetIsAdmin(true).
			SetRole(user.RoleAdmin).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
	} else if err == nil && (!adminUser.IsAdmin || adminUser.Role != user.RoleAdmin) {
		adminUser, _ = adminUser.Update().SetIsAdmin(true).SetRole(user.RoleAdmin).Save(ctx)
	}

	hashedUserPassword, _ := bcrypt.GenerateFromPassword([]byte("User@123456"), bcrypt.DefaultCost)
	regularUser, err := client.User.Query().Where(user.EmailEQ("user@fooddiscovery.com")).Only(ctx)
	if ent.IsNotFound(err) {
		regularUser, err = client.User.Create().
			SetName("Regular Foodie").
			SetEmail("user@fooddiscovery.com").
			SetPassword(string(hashedUserPassword)).
			SetVerified(true).
			SetIsAdmin(false).
			SetRole(user.RoleCustomer).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create regular user: %w", err)
		}
	} else if err == nil && regularUser.Role != user.RoleCustomer {
		regularUser, _ = regularUser.Update().SetRole(user.RoleCustomer).Save(ctx)
	}

	hashedMerchantPassword, _ := bcrypt.GenerateFromPassword([]byte("Merchant@123456"), bcrypt.DefaultCost)
	merchantUser, err := client.User.Query().Where(user.EmailEQ("merchant.ctr@fooddiscovery.com")).Only(ctx)
	if ent.IsNotFound(err) {
		merchantUser, err = client.User.Create().
			SetName("CTR Malleshwaram Owner").
			SetEmail("merchant.ctr@fooddiscovery.com").
			SetPassword(string(hashedMerchantPassword)).
			SetVerified(true).
			SetIsAdmin(false).
			SetRole(user.RoleMerchant).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create merchant user: %w", err)
		}
	} else if err == nil && (merchantUser.Role != user.RoleMerchant || merchantUser.IsAdmin) {
		merchantUser, _ = merchantUser.Update().SetRole(user.RoleMerchant).SetIsAdmin(false).Save(ctx)
	}


	// 2. Seed City (Bangalore)
	blr, err := client.City.Query().Where(city.NameEQ("Bangalore")).Only(ctx)
	if ent.IsNotFound(err) {
		blr, err = client.City.Create().
			SetName("Bangalore").
			SetState("Karnataka").
			SetLatitude(12.9716).
			SetLongitude(77.5946).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create Bangalore city: %w", err)
		}
	}

	// Helper to check and create restaurant with food items
	seedResto := func(
		name, description, address string,
		lat, lng float64,
		phone, hours string,
		status restaurant.Status,
		foods []struct {
			name     string
			desc     string
			category string
			price    float64
			avail    bool
		},
	) error {
		existing, err := client.Restaurant.Query().Where(restaurant.NameEQ(name)).Only(ctx)
		if ent.IsNotFound(err) {
			r, err := client.Restaurant.Create().
				SetName(name).
				SetDescription(description).
				SetAddress(address).
				SetLatitude(lat).
				SetLongitude(lng).
				SetPhone(phone).
				SetOpeningHours(hours).
				SetStatus(status).
				SetCity(blr).
				Save(ctx)
			if err != nil {
				return err
			}

			for _, f := range foods {
				_, err = client.FoodItem.Create().
					SetName(f.name).
					SetDescription(f.desc).
					SetCategory(f.category).
					SetPrice(f.price).
					SetAvailable(f.avail).
					SetRestaurant(r).
					Save(ctx)
				if err != nil {
					return err
				}
			}
		} else if err == nil && existing != nil {
			// Update status in case it was updated
			_, _ = existing.Update().SetStatus(status).Save(ctx)
		}
		return nil
	}

	// 3. Seed sample Bangalore restaurants
	// Restaurant 1: Vidyarthi Bhavan (Gandhi Bazaar) - Has ₹100 items (Approved)
	_ = seedResto(
		"Vidyarthi Bhavan",
		"Legendary heritage eatery famous for thick, crispy Benne Masala Dosa and filter coffee.",
		"32, Gandhi Bazaar Main Rd, Basavanagudi, Bengaluru, Karnataka 560004",
		12.9451, 77.5714,
		"+91 80 2667 7588", "6:30 AM - 11:30 AM, 2:00 PM - 8:00 PM",
		restaurant.StatusApproved,
		[]struct {
			name     string
			desc     string
			category string
			price    float64
			avail    bool
		}{
			{"Crispy Masala Dosa", "Signature ghee-roasted dosa with potato filling", "South Indian", 70.0, true},
			{"Idli Vada Combo", "Steamed rice cakes with golden crispy vada and coconut chutney", "Breakfast", 55.0, true},
			{"Poori Saagu", "Puffy pooris with mildly spiced mixed vegetable gravy", "Breakfast", 60.0, true},
			{"Khara Bath", "Spiced semolina upma tempered with curry leaves and mustard", "Breakfast", 40.0, true},
			{"Special Filter Coffee", "Authentic South Indian chicory-blend filter kaapi", "Beverages", 25.0, true},
			{"Kesari Bath", "Sweet semolina dessert with saffron, ghee, and roasted cashews", "Dessert", 45.0, true},
		},
	)

	// Restaurant 2: CTR (Shri Sagar) - Malleshwaram - Has ₹100 items (Approved)
	_ = seedResto(
		"Central Tiffin Room (CTR) - Shri Sagar",
		"Iconic North Bangalore breakfast spot acclaimed for butter-laden benne masala dosas.",
		"7th Cross, Margosa Rd, Malleshwaram, Bengaluru, Karnataka 560003",
		12.9982, 77.5703,
		"+91 80 2331 7531", "7:30 AM - 12:30 PM, 4:00 PM - 9:00 PM",
		restaurant.StatusApproved,
		[]struct {
			name     string
			desc     string
			category string
			price    float64
			avail    bool
		}{
			{"Benne Masala Dosa", "Crispy exterior, soft inside with dollops of fresh butter", "South Indian", 75.0, true},
			{"Mangalore Bajji (4 pcs)", "Crispy golden fried maida fritters served with green chutney", "Snacks", 45.0, true},
			{"Filter Coffee", "Freshly brewed strong aromatic filter coffee", "Beverages", 25.0, true},
			{"Rava Idli", "Steamed semolina dumpling with cashew and coriander", "South Indian", 45.0, true},
			{"Chow Chow Bath", "Duo of savory Khara Bath and sweet Kesari Bath", "Combo", 80.0, true},
		},
	)

	// Restaurant 3: Brahmin's Coffee Bar - Shankarapuram - Has ₹100 items (Approved)
	_ = seedResto(
		"Brahmin's Coffee Bar",
		"Beloved standing breakfast haven serving heavenly soft idlis and unbeatable coconut chutney.",
		"Ranga Rao Rd, near Shankar Mutt, Shankarapuram, Bengaluru, Karnataka 560004",
		12.9542, 77.5689,
		"+91 98450 30234", "6:00 AM - 12:00 PM, 3:00 PM - 7:00 PM",
		restaurant.StatusApproved,
		[]struct {
			name     string
			desc     string
			category string
			price    float64
			avail    bool
		}{
			{"Button Idli (2 pcs)", "Melt-in-mouth steamed idlis smothered in fresh coconut chutney", "Breakfast", 35.0, true},
			{"Crispy Medu Vada", "Hot, crunchy urad dal vada", "Breakfast", 30.0, true},
			{"Khara Bath", "Traditional Karnataka style upma with vegetables", "Breakfast", 40.0, true},
			{"Kesari Bath", "Ghee-rich pineapple semolina sweet", "Dessert", 40.0, true},
			{"Filter Kaapi", "Piping hot frothy filter coffee in steel tumbler", "Beverages", 20.0, true},
		},
	)

	// Restaurant 4: Taaza Thindi - Jayanagar - Has ultra-budget items <= ₹100 (Approved)
	_ = seedResto(
		"Taaza Thindi",
		"Hyper-clean modern quick-service outlet offering premium South Indian fare at budget prices.",
		"1004, 26th Main Rd, 4th T Block, Jayanagar, Bengaluru, Karnataka 560041",
		12.9279, 77.5834,
		"+91 80 4091 2277", "7:00 AM - 12:00 PM, 4:30 PM - 9:30 PM",
		restaurant.StatusApproved,
		[]struct {
			name     string
			desc     string
			category string
			price    float64
			avail    bool
		}{
			{"Masala Dosa", "Golden brown crisp dosa with potato palya and mint chutney", "South Indian", 35.0, true},
			{"Plain Dosa", "Thin crisp crepe with chutney", "South Indian", 30.0, true},
			{"Steamed Idli (2 pcs)", "Ultra soft fluffy idlis with coconut chutney", "Breakfast", 20.0, true},
			{"Crispy Vada", "Deep fried spiced lentil doughnut", "Breakfast", 20.0, true},
			{"Chow Chow Bath", "Khara Bath + Kesari Bath combination plate", "Breakfast", 40.0, true},
			{"Pure Filter Coffee", "Aromatic South Indian filter coffee", "Beverages", 15.0, true},
		},
	)

	// Restaurant 5: Hari Super Sandwich - Jayanagar - Chaat & Street Food <= ₹100 (Approved)
	_ = seedResto(
		"Hari Super Sandwich",
		"Bangalore's cult street sandwich and chaat destination known for innovative toasts.",
		"177/A-44, 22nd Cross Rd, 3rd Block, Jayanagar, Bengaluru, Karnataka 560011",
		12.9260, 77.5840,
		"+91 99011 11330", "11:00 AM - 10:00 PM",
		restaurant.StatusApproved,
		[]struct {
			name     string
			desc     string
			category string
			price    float64
			avail    bool
		}{
			{"Chocolate Sandwich", "Toasted sandwich loaded with melted dairy milk chocolate", "Sandwiches", 70.0, true},
			{"Cheese Chilli Toast", "Crispy toast topped with green chillies and grated Amul cheese", "Sandwiches", 85.0, true},
			{"Dahi Puri", "Crisp puris filled with spiced potatoes, sweet yogurt, and sev", "Chaat", 60.0, true},
			{"Pani Puri", "Crispy hollow puris served with tangy mint water and potato mix", "Chaat", 45.0, true},
			{"Sev Puri", "Flat puris topped with onions, tomatoes, chutneys, and nylon sev", "Chaat", 50.0, true},
			{"American Sweet Corn Sandwich", "Sweet corn and spiced mayo grilled sandwich", "Sandwiches", 75.0, true},
		},
	)

	// Restaurant 6: Empire Restaurant - Church Street - Quick bites & Rolls <= ₹100 (Approved)
	_ = seedResto(
		"Empire Restaurant - Church Street",
		"Popular Bangalore late-night diner with quick kebab rolls and street delicacies.",
		"36, Church St, Shanthala Nagar, Ashok Nagar, Bengaluru, Karnataka 560001",
		12.9750, 77.6033,
		"+91 80 4041 4041", "11:00 AM - 1:00 AM",
		restaurant.StatusApproved,
		[]struct {
			name     string
			desc     string
			category string
			price    float64
			avail    bool
		}{
			{"Chicken Shawarma Roll", "Shredded grilled chicken in rumali roti with garlic mayo", "Rolls", 99.0, true},
			{"Paneer Tikka Roll", "Charcoal grilled cottage cheese chunks with mint sauce in flatbread", "Rolls", 90.0, true},
			{"Egg Kothu Parotta", "Minced flaky parotta tossed with eggs and spicy salna", "Street Food", 95.0, true},
			{"Empire Special Chicken Biryani", "Fragrant spiced chicken biryani with boiled egg and raita", "Main Course", 290.0, true},
			{"Butter Chicken", "Rich creamy tomato gravy with tender chicken pieces", "Curries", 320.0, true},
		},
	)

	// Restaurant 7: Corner House Ice Cream - Indiranagar - Desserts with items <= ₹100 (Approved)
	_ = seedResto(
		"Corner House Ice Cream - Indiranagar",
		"Bengaluru's foremost ice cream parlour brand beloved across generations.",
		"470, 80 Feet Rd, HAL 3rd Stage, Indiranagar, Bengaluru, Karnataka 560075",
		12.9784, 77.6408,
		"+91 98450 11222", "11:00 AM - 11:30 PM",
		restaurant.StatusApproved,
		[]struct {
			name     string
			desc     string
			category string
			price    float64
			avail    bool
		}{
			{"Single Scoop Vanilla Sundae", "Pure dairy vanilla scoop with warm hot chocolate drizzle", "Ice Cream", 70.0, true},
			{"Mini Hot Chocolate Fudge", "Classic HCF portion with peanuts and fudge sauce", "Sundae", 95.0, true},
			{"Mango Mousse Cup", "Seasonal mango cream dessert in a cup", "Dessert", 80.0, true},
			{"Death By Chocolate (DBC)", "Full size rich chocolate cake, vanilla ice cream, cherries & cream", "Sundae", 240.0, true},
			{"Apple Pie with Ice Cream", "Warm spiced apple pie with single vanilla scoop", "Dessert", 160.0, true},
		},
	)

	// Restaurant 8: Truffles - Koramangala - Appetizers <= ₹100 (Approved)
	_ = seedResto(
		"Truffles - Koramangala",
		"High-energy cafe renowned for gourmet burgers, pastas, and savory snacks.",
		"93, Ground Floor, 4th B Cross, 5th Block, Koramangala, Bengaluru, Karnataka 560095",
		12.9345, 77.6200,
		"+91 80 4146 6565", "11:00 AM - 11:00 PM",
		restaurant.StatusApproved,
		[]struct {
			name     string
			desc     string
			category string
			price    float64
			avail    bool
		}{
			{"Crisp Potato Wedges", "Seasoned golden potato wedges with garlic dip", "Starters", 95.0, true},
			{"Cheesy Garlic Bread", "French loaf brushed with herb butter and melted mozzarella", "Starters", 99.0, true},
			{"All American Cheese Burger", "Juicy grilled patty with cheddar cheese and gherkins", "Burgers", 220.0, true},
			{"Ferrero Rocher Milkshake", "Decadent thick hazelnut chocolate shake", "Beverages", 190.0, true},
		},
	)
	// 9. Veena Stores (Malleshwaram)
	if err := seedResto(
			"Veena Stores",
			"Famous street-corner tiffin spot running since 1977, known for fluffy idlis, crispy vadas, and fresh mint-coconut chutney.",
			"15th Cross, Margosa Rd, Malleshwaram, Bengaluru, Karnataka 560003",
			13.0062, 77.5702,
			"+91 80 2334 4838",
			"6:30 AM - 12:00 PM, 3:30 PM - 9:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Steamed Button Idli (2 pcs)", "Soft fluffy rice cakes with fresh coconut-mint chutney", "Breakfast", 25.0, true},
				{"Crispy Uddin Vada", "Deep fried spiced urad dal fritter", "Breakfast", 30.0, true},
				{"Khara Bath (Upma)", "Spicy semolina tempered with mustard, curry leaves, and veggies", "Breakfast", 30.0, true},
				{"Sweet Kesari Bath", "Saffron pineapple semolina halwa with pure ghee", "Dessert", 30.0, true},
				{"Shavige Bath", "Lemon spiced rice vermicelli with grated coconut", "Breakfast", 35.0, true},
				{"Authentic Filter Kaapi", "Traditional frothy South Indian chicory blend", "Beverages", 15.0, true},
			},
		); err != nil {
			return err
		}

		// 10. Asha Tiffins (HSR Layout)
		if err := seedResto(
			"Asha Tiffins - HSR Layout",
			"Modern Darshini sensation in HSR Layout beloved for podi-slathered dosas, button idlis, and lightning-fast service.",
			"5th Sector, 14th Main Rd, HSR Layout, Bengaluru, Karnataka 560102",
			12.9152, 77.6385,
			"+91 80 4370 7888",
			"7:00 AM - 10:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Ghee Podi Masala Dosa", "Golden crisp dosa coated in aromatic gun powder podi and butter", "South Indian", 75.0, true},
				{"Ghee Button Idlis (14 pcs)", "Mini idlis immersed in piping hot sambar with ghee", "Breakfast", 45.0, true},
				{"Crispy Medu Vada", "Golden brown lentil doughnut with coconut chutney", "Breakfast", 30.0, true},
				{"Poori Saagu (3 pcs)", "Puffed whole wheat bread with spiced potato vegetable curry", "Breakfast", 60.0, true},
				{"Filter Kaapi", "Aromatic South Indian filter coffee in steel tumbler", "Beverages", 15.0, true},
				{"Sweet Kesari Bath", "Pure ghee pineapple halwa", "Dessert", 35.0, true},
			},
		); err != nil {
			return err
		}

		// 11. Indian Coffee House (Church Street)
		if err := seedResto(
			"Indian Coffee House - Church Street",
			"Vintage heritage cafe operating since 1957 with colonial charm, waiters in Nehru caps, and budget cutlets.",
			"57, Church St, Shanthala Nagar, Ashok Nagar, Bengaluru, Karnataka 560001",
			12.9748, 77.6045,
			"+91 80 2558 7088",
			"8:30 AM - 8:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Mutton Cutlet (2 pcs)", "Crispy spiced minced meat croquettes with tomato sauce", "Snacks", 65.0, true},
				{"Vegetable Cutlet", "Golden breaded vegetable patties with beetroot tinge", "Snacks", 45.0, true},
				{"Masala Cheese Omelette", "Fluffy two-egg omelette with onion, green chillies, and toast", "Breakfast", 60.0, true},
				{"Cold Coffee with Cream", "Thick chilled brewed coffee shake", "Beverages", 50.0, true},
				{"Heritage Hot Filter Coffee", "Classic brew served in retro crockery", "Beverages", 25.0, true},
			},
		); err != nil {
			return err
		}

		// 12. Albert Bakery (Mosque Road, Frazer Town)
		if err := seedResto(
			"Albert Bakery - Mosque Road",
			"Century-old bakery operating since 1902, famed across Bangalore for mutton keema samosas and fresh hot puffs.",
			"93, Mosque Rd, Cleveland Town, Pulikeshi Nagar, Bengaluru, Karnataka 560005",
			12.9989, 77.6162,
			"+91 98861 65349",
			"3:00 PM - 9:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Mutton Keema Samosa (2 pcs)", "Famous triangular pastry stuffed with spicy minced mutton", "Snacks", 35.0, true},
				{"Mutton Brain Puff", "Crispy flaky puff with delicately seasoned bheja fry", "Snacks", 60.0, true},
				{"Egg Puff", "Flaky puff pastry enclosing spiced hard-boiled egg", "Snacks", 30.0, true},
				{"Khova Naan", "Rich sweet baked pastry with condensed milk filling", "Dessert", 50.0, true},
				{"Coconut Pastry", "Fresh baked coconut cake slice", "Dessert", 35.0, true},
			},
		); err != nil {
			return err
		}

		// 13. Cool Joint (Jayanagar)
		if err := seedResto(
			"Cool Joint - Jayanagar",
			"Beloved 4th Block Jayanagar landmark famous for Bun Gulkand, fresh fruit sundaes, and grilled cheese toasts.",
			"11th Main Rd, 4th Block, Jayanagar, Bengaluru, Karnataka 560011",
			12.9288, 77.5832,
			"+91 80 2663 3224",
			"9:00 AM - 10:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Bun Gulkand with Ice Cream", "Sweet rose petal preserve in butter bun topped with vanilla ice cream", "Dessert", 60.0, true},
				{"Fresh Fruit Salad with Ice Cream", "Diced seasonal fruits with vanilla scoop", "Dessert", 70.0, true},
				{"Spicy Masala Cheese Toast", "Grilled sandwich with capsicum, onion, and melted cheese", "Sandwiches", 55.0, true},
				{"Fresh Mosambi Sweet Lime Juice", "Freshly squeezed citrus juice", "Beverages", 50.0, true},
			},
		); err != nil {
			return err
		}

		// 14. Bangalore Iyengar Bakery
		if err := seedResto(
			"Bangalore Iyengar Bakery",
			"Historic bakery tradition dating to 1898, famous for warm vegetable puffs, khara buns, and honey cake slices.",
			"Sampige Rd, Malleshwaram, Bengaluru, Karnataka 560003",
			12.9958, 77.5715,
			"+91 80 2331 1290",
			"7:00 AM - 10:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Classic Vegetable Puff", "Spicy onion and potato stuffed layered puff", "Snacks", 20.0, true},
				{"Spicy Potato Bun (Aloo Bun)", "Soft bun stuffed with dry potato masala", "Snacks", 25.0, true},
				{"Khara Butter Bun", "Spiced masala bun slathered with pure salted butter", "Snacks", 20.0, true},
				{"Sweet Dilkush Pastry", "Baked bun loaded with tutti frutti and shredded coconut", "Dessert", 35.0, true},
				{"Honey Cake Slice", "Moist sponge cake with honey glaze and coconut flakes", "Dessert", 25.0, true},
			},
		); err != nil {
			return err
		}

		// 15. By2 Coffee (Basavanagudi)
		if err := seedResto(
			"By2 Coffee - Basavanagudi",
			"Traditional South Indian coffee joint celebrating the culture of sharing a cup of kaapi and hot idlis.",
			"Bugle Rock Rd, Basavanagudi, Bengaluru, Karnataka 560004",
			12.9422, 77.5684,
			"+91 80 2660 1199",
			"6:30 AM - 12:30 PM, 3:30 PM - 9:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"By2 Filter Coffee (Half)", "Strong frothy decoction in traditional tumbler", "Beverages", 15.0, true},
				{"Mini Idli with Chutney", "Two soft steamed cakes with hot spicy chutney", "Breakfast", 20.0, true},
				{"Crispy Vada", "Hot golden lentil doughnut", "Breakfast", 25.0, true},
				{"Chow Chow Bath", "Twin combo of savory Khara Bath and sweet Kesari Bath", "Breakfast", 40.0, true},
			},
		); err != nil {
			return err
		}

		// 16. Sendhoor Coffee (Indiranagar)
		if err := seedResto(
			"Sendhoor Coffee - Indiranagar",
			"Tamil Nadu style tiffin bar popular for paper-crisp ghee roast dosas, podi idlis, and spicy tomato kurma.",
			"Double Rd, Stage 2, Indiranagar, Bengaluru, Karnataka 560038",
			12.9712, 77.6415,
			"+91 97424 55333",
			"6:30 AM - 10:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Tamil Nadu Ghee Roast Dosa", "Thin paper-crisp conical dosa roasted in pure ghee", "South Indian", 65.0, true},
				{"Mini Ghee Podi Idli (10 pcs)", "Tossed in fiery spicy molagapodi and clarified butter", "Breakfast", 40.0, true},
				{"Ven Pongal with Chutney & Sambar", "Peppercorn and cumin tempered moong dal rice with ghee", "Breakfast", 50.0, true},
				{"Filter Kaapi", "Authentic Kumbakonam degree-style coffee", "Beverages", 15.0, true},
				{"Parotta Kurma (2 pcs)", "Flaky layered parotta with spicy vegetable coconut gravy", "South Indian", 60.0, true},
			},
		); err != nil {
			return err
		}

		// 17. Sri Sairam's Chats (Malleshwaram)
		if err := seedResto(
			"Sri Sairam's Chats and Juice",
			"Inventive street chaat institution famous for signature wafer-based chips masala and disco chaats.",
			"16th Cross Rd, Margosa Rd, Malleshwaram, Bengaluru, Karnataka 560003",
			13.0075, 77.5698,
			"+91 80 2334 0929",
			"2:30 PM - 10:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Crazy Masala Chaat", "Crispy chips, crushed papdi, spicy tangy gravy and raw onions", "Chaat", 50.0, true},
				{"Chips Masala", "Potato wafers topped with spicy warm peas ragda and sev", "Chaat", 45.0, true},
				{"Dahi Batata Puri", "Puffed puris stuffed with potatoes and sweetened chilled curd", "Chaat", 50.0, true},
				{"Masala Puri", "Crushed puris drenched in boiling spicy green peas curry", "Chaat", 40.0, true},
				{"Fresh Mango Milkshake", "Chilled Alphonso pulp blended with milk", "Beverages", 60.0, true},
			},
		); err != nil {
			return err
		}

		// 18. Mahalakshmi Tiffin Room (Gandhi Bazaar)
		if err := seedResto(
			"Mahalakshmi Tiffin Room (MLTR)",
			"Historic 1926 Gandhi Bazaar eatery famous for classic wood-roasted masala dosas and rich curd vadas.",
			"85, DVG Rd, Gandhi Bazaar, Basavanagudi, Bengaluru, Karnataka 560004",
			12.9463, 77.5721,
			"+91 80 2667 9231",
			"7:00 AM - 12:30 PM, 4:00 PM - 8:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Crispy Masala Dosa", "Traditional thick golden crust with potato palya", "South Indian", 65.0, true},
				{"Curd Vada (1 pc)", "Vada soaked in thick spiced yogurt tempered with mustard and boondi", "Snacks", 45.0, true},
				{"Rava Idli with Sagu", "Steamed semolina cake with cashew nuts and vegetable curry", "Breakfast", 45.0, true},
				{"Filter Coffee", "Hot traditional filter brew", "Beverages", 25.0, true},
			},
		); err != nil {
			return err
		}

		// 19. Janatha Hotel (Malleshwaram)
		if err := seedResto(
			"Janatha Hotel - Malleshwaram",
			"Old-school dining hall on 8th Cross serving crispy red-chutney masala dosas and hot filter coffee.",
			"8th Cross Rd, Malleshwaram, Bengaluru, Karnataka 560003",
			13.0001, 77.5708,
			"+91 80 2334 5609",
			"7:30 AM - 1:00 PM, 3:30 PM - 8:45 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Golden Masala Dosa", "Crispy red-chilli spiced butter roast", "South Indian", 60.0, true},
				{"Rava Dosa", "Lacy crispy semolina crepe with coconut chutney", "South Indian", 70.0, true},
				{"Medu Vada Sambar", "Crispy lentil fritter dunked in hot sambar", "Breakfast", 35.0, true},
				{"Chicory Filter Coffee", "Frothy heritage brew", "Beverages", 20.0, true},
			},
		); err != nil {
			return err
		}

		// 20. New Krishna Bhavan (Malleshwaram)
		if err := seedResto(
			"New Krishna Bhavan",
			"Renowned coastal Karnataka and Udupi specialty restaurant serving Mangalore buns and green masala dosas.",
			"Opposite Mantri Mall, Sampige Rd, Malleshwaram, Bengaluru, Karnataka 560003",
			12.9912, 77.5722,
			"+91 80 2344 6500",
			"7:30 AM - 9:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Button Idli Sambar (12 pcs)", "Mini idlis floating in aromatic Udupi sambar", "Breakfast", 55.0, true},
				{"Green Masala Dosa", "Spiced with fresh mint and coriander chutney inside", "South Indian", 75.0, true},
				{"Mangalore Buns (2 pcs)", "Sweet banana puffed bread served with spicy coconut chutney", "Breakfast", 45.0, true},
				{"Bisibelebath with Boondi", "Spiced rice and lentil mash with pure ghee and crispy pearls", "South Indian", 60.0, true},
			},
		); err != nil {
			return err
		}

		// 21. Al-Amanah Cafe (Koramangala)
		if err := seedResto(
			"Al-Amanah Cafe - Koramangala",
			"Cult-favorite shawarma haven in 5th Block Koramangala known for jumbo chicken rolls and spiced fries.",
			"5th Block, KHB Colony, Koramangala, Bengaluru, Karnataka 560095",
			12.9348, 77.6189,
			"+91 80 4110 5005",
			"11:00 AM - 11:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Khubz Chicken Shawarma Roll", "Shredded grilled chicken with garlic toum in soft pita", "Rolls", 95.0, true},
				{"Falafel Roll with Hummus", "Crispy chickpea fritters with pickled veggies and tahini", "Rolls", 70.0, true},
				{"Peri Peri Salted French Fries", "Crispy golden potato fries with spicy seasoning", "Fast Food", 70.0, true},
				{"Fresh Mint Lime Cooler", "Chilled refreshing lime soda with fresh mint", "Beverages", 40.0, true},
			},
		); err != nil {
			return err
		}

		// 22. VV Puram Food Street - VB Bakery
		if err := seedResto(
			"VB Bakery - VV Puram",
			"The culinary anchor of Bangalore's famous Food Street, operating since 1953 and renowned for Dumrot Halwa.",
			"Sajjan Rao Circle, VV Puram, Bengaluru, Karnataka 560004",
			12.9515, 77.5772,
			"+91 80 2667 7055",
			"6:00 AM - 10:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Dumrot Halwa (Ash Gourd)", "Rich caramel ash gourd halwa dripping with pure ghee", "Dessert", 40.0, true},
				{"Congress Kadlekai Bun", "Sweet soft bun stuffed with spiced split peanuts and butter", "Snacks", 35.0, true},
				{"Khara Butter Biscuit (pack)", "Spicy savory cumin butter biscuits", "Snacks", 30.0, true},
			},
		); err != nil {
			return err
		}

		// 23. VV Puram Food Street - Ramu Butter Dosa
		if err := seedResto(
			"Ramu Butter Dosa - Food Street",
			"The legendary butter dosa cart on Thindi Beedi (Food Street) known for giant dollops of butter on sizzling roasts.",
			"Old Hospital Rd, Food Street, VV Puram, Bengaluru, Karnataka 560004",
			12.9519, 77.5778,
			"+91 98451 22345",
			"5:30 PM - 11:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Butter Ghee Roast Dosa", "Smothered in a giant block of fresh melting butter", "South Indian", 70.0, true},
				{"Red Chilli Pudi Dosa", "Loaded with spicy gun powder and potato filling", "South Indian", 65.0, true},
				{"Hot Paddus / Guliyappa (6 pcs)", "Crispy spherical rice dumplings with coconut chutney", "Snacks", 50.0, true},
			},
		); err != nil {
			return err
		}

		// 24. Sri Udupi Park (Whitefield)
		if err := seedResto(
			"Sri Udupi Park - Whitefield",
			"The premier vegetarian Darshini hub on ITPL Main Road serving quick, hygienic hot meals all day.",
			"ITPL Main Rd, Kundalahalli Colony, Brookefield, Bengaluru, Karnataka 560037",
			12.9702, 77.7145,
			"+91 80 4209 8899",
			"6:30 AM - 11:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Crispy Masala Dosa", "Classic butter roast with potato mash", "South Indian", 65.0, true},
				{"Soft Set Dosa (3 pcs)", "Spongey soft dosas with vegetable saagu and chutney", "South Indian", 55.0, true},
				{"Steamed Rava Idli", "Semolina cake with cashew nuts and ghee", "Breakfast", 40.0, true},
				{"Bisibelebath with Khara Boondi", "Hot spicy lentil rice dish", "South Indian", 50.0, true},
				{"Filter Coffee", "Hot frothy South Indian coffee", "Beverages", 20.0, true},
			},
		); err != nil {
			return err
		}

		// 25. Thoms Bakery (Frazer Town)
		if err := seedResto(
			"Thoms Bakery - Frazer Town",
			"Iconic heritage bakery on Wheeler Road famous for warm savory mutton puffs and traditional plum cakes.",
			"1/2, Wheeler Rd, Cleveland Town, Frazer Town, Bengaluru, Karnataka 560005",
			12.9972, 77.6185,
			"+91 80 2536 1076",
			"8:30 AM - 9:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Fresh Mutton Puff", "Flaky puff pastry stuffed with spicy minced meat", "Snacks", 55.0, true},
				{"Chicken Puff", "Layered pastry filled with pepper shredded chicken", "Snacks", 45.0, true},
				{"Egg Puff", "Classic egg filling with caramelized onions", "Snacks", 30.0, true},
				{"Plum Cake Slice", "Traditional rich spiced plum fruit cake", "Dessert", 40.0, true},
			},
		); err != nil {
			return err
		}

		// 26. Kolkata Famous Kati Roll (Church Street)
		if err := seedResto(
			"Kolkata Famous Kati Roll",
			"Street-style parotta roll counter on Church Street serving crisp egg and chicken rolls with zesty sauces.",
			"12, Church St, Shanthala Nagar, Ashok Nagar, Bengaluru, Karnataka 560001",
			12.9752, 77.6052,
			"+91 98311 23456",
			"11:00 AM - 11:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Double Egg Roll", "Two farm-fresh eggs on layered parotta with tangy onions", "Rolls", 60.0, true},
				{"Single Chicken Roll", "Tender marinated chicken chunks wrapped with lime and chat masala", "Rolls", 80.0, true},
				{"Paneer Tikka Roll", "Spiced grilled cottage cheese with peppers in parotta", "Rolls", 80.0, true},
				{"Aloo Roll", "Spicy mashed potato roll with green chillies", "Rolls", 60.0, true},
			},
		); err != nil {
			return err
		}

		// 27. Rakesh Kumar Pani Puri (Jayanagar)
		if err := seedResto(
			"Rakesh Kumar Pani Puri - Jayanagar",
			"High-energy street cart famous for authentic Kolkata-style crisp puchkas and spicy churmur chaat.",
			"8th F Main Rd, 3rd Block, Jayanagar, Bengaluru, Karnataka 560011",
			12.9272, 77.5845,
			"+91 98450 11223",
			"2:00 PM - 10:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Kolkata Style Puchka (6 pcs)", "Crisp hollow wheat shells with spiced potato mash and mint tamarind water", "Chaat", 40.0, true},
				{"Spicy Churmur Chaat", "Crushed puchkas with boiled potatoes, mustard oil, and spices", "Chaat", 50.0, true},
				{"Dahi Papdi Chaat", "Crispy flat flour wafers topped with yogurt, tamarind chutney, and sev", "Chaat", 55.0, true},
			},
		); err != nil {
			return err
		}

		// 28. Fanoos / Siddique Kabab Centre
		if err := seedResto(
			"Siddique Kabab Centre - Johnson Market",
			"Legendary Richmond Town eatery famous for charcoal-grilled sheekh rolls and piping hot beef/chicken samosas.",
			"Johnson Market, Hosur Rd, Richmond Town, Bengaluru, Karnataka 560025",
			12.9615, 77.6082,
			"+91 80 2222 5556",
			"12:00 PM - 11:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Chicken Sheekh Roll", "Minced spiced chicken kebab grilled on charcoal in soft rumali", "Rolls", 85.0, true},
				{"Mutton Samosa (2 pcs)", "Crispy fried patti samosa stuffed with spicy minced meat", "Snacks", 30.0, true},
				{"Chicken Shawarma in Rumali", "Spit-roasted chicken with garlic cream and pickles", "Rolls", 90.0, true},
			},
		); err != nil {
			return err
		}

		// 29. Chikkanna Tiffin Room
		if err := seedResto(
			"Chikkanna Tiffin Room - City Market",
			"Iconic heritage spot near City Market famous for soft, spongey ghee khali dosas and fiery podi idlis.",
			"Nagarthpet Main Rd, City Market, Bengaluru, Karnataka 560002",
			12.9658, 77.5812,
			"+91 80 2223 4455",
			"7:00 AM - 12:30 PM, 4:30 PM - 8:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Ghee Khali Dosa (2 pcs)", "Ultra-soft spongey fermented rice crepes with hot sagu and chutney", "South Indian", 60.0, true},
				{"Mini Ghee Masala Dosa", "Crispy butter dosa stuffed with spicy onion palya", "South Indian", 65.0, true},
				{"Spicy Podi Idli", "Tossed in fiery gun powder and hot melted ghee", "Breakfast", 45.0, true},
			},
		); err != nil {
			return err
		}

		// 30. Tibetan Mother's Kitchen (Koramangala)
		if err := seedResto(
			"Tibetan Mother's Kitchen - Koramangala",
			"Authentic Himalayan cozy food joint famous for steaming momos with fiery red chili sauce.",
			"1st Main Rd, 7th Block, Koramangala, Bengaluru, Karnataka 560095",
			12.9362, 77.6152,
			"+91 98860 77889",
			"12:00 PM - 10:30 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Steamed Veg Momos (8 pcs)", "Minced cabbage and vegetables in thin wrapper with spicy red chili dip", "Fast Food", 70.0, true},
				{"Steamed Chicken Momos (8 pcs)", "Juicy minced chicken dumplings with hot sesame dip", "Fast Food", 90.0, true},
				{"Veg Thukpa Noodle Soup (Half)", "Himalayan noodle soup with ginger, garlic, and fresh vegetables", "Fast Food", 95.0, true},
			},
		); err != nil {
			return err
		}

		// 31. Bhaskar's Mane Holige
		if err := seedResto(
			"Bhaskar's Mane Holige - VV Puram",
			"Specialty sweet shop offering hot, live-prepared coconut and dal holiges (sweet obbattu) dripping with ghee.",
			"Sajjan Rao Circle, VV Puram, Bengaluru, Karnataka 560004",
			12.9511, 77.5775,
			"+91 80 2661 2233",
			"10:00 AM - 10:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Hot Coconut Holige (Obbattu)", "Traditional sweet flatbread stuffed with jaggery and fresh coconut", "Dessert", 25.0, true},
				{"Bili Bele Dal Holige", "Gram dal and cardamom stuffed sweet bread served warm with ghee", "Dessert", 25.0, true},
				{"Dry Fruit Special Holige", "Enriched with cashews, almonds, and saffron", "Dessert", 40.0, true},
			},
		); err != nil {
			return err
		}

		// 32. The Old Fashioned Chai Co. (Indiranagar)
		if err := seedResto(
			"The Old Fashioned Chai Co.",
			"Cozy tea and quick-snack corner on 100ft road serving fresh cutting adrak chai and warm bun maska.",
			"100 Feet Rd, HAL 2nd Stage, Indiranagar, Bengaluru, Karnataka 560038",
			12.9735, 77.6410,
			"+91 80 4123 9988",
			"8:00 AM - 11:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Adrak Ginger Cutting Chai", "Freshly brewed milk tea infused with crushed ginger", "Beverages", 25.0, true},
				{"Cardamom Masala Chai", "Aromatic spiced tea with cloves and cardamom", "Beverages", 25.0, true},
				{"Warm Bun Maska", "Soft sweet bun sliced and slathered with pure amul butter", "Snacks", 35.0, true},
				{"Classic Butter Maggi", "Noodles cooked with spices, butter, and green peas", "Snacks", 50.0, true},
			},
		); err != nil {
			return err
		}

		// 33. Shivanna Gulkand Center (VV Puram) -> PENDING APPROVAL FOR ADMIN
		if err := seedResto(
			"Shivanna Gulkand Center",
			"Iconic Food Street dessert shop serving rose petal gulkand with butter, dry fruits, and ice cream.",
			"Near Sajjan Rao Circle, Food Street, VV Puram, Bengaluru, Karnataka 560004",
			12.9518, 77.5774,
			"+91 98452 33445",
			"5:00 PM - 11:30 PM",
			restaurant.StatusPending,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Special Gulkand with Vanilla Scoop", "Pure rose petal preserve with dry fruits, fig, and ice cream", "Dessert", 70.0, true},
				{"Fresh Fruit Salad with Gulkand", "Seasonal cut fruits with rose preserve", "Dessert", 65.0, true},
			},
		); err != nil {
			return err
		}

		// 34. AECS Layout Roll Corner (Whitefield) -> PENDING APPROVAL FOR ADMIN
		if err := seedResto(
			"AECS Layout Roll Corner",
			"Popular tech-corridor food truck serving hot Kolkata-style egg and chicken rolls.",
			"AECS Layout, Brookefield, Bengaluru, Karnataka 560037",
			12.9691, 77.7160,
			"+91 97412 88990",
			"4:00 PM - 11:00 PM",
			restaurant.StatusPending,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Egg Chicken Kathi Roll", "Spiced chicken cubes wrapped in egg-fried parotta", "Rolls", 90.0, true},
				{"Paneer Tikka Roll", "Charred paneer cubes with onion and mint spread", "Rolls", 80.0, true},
			},
		); err != nil {
			return err
		}

		// 35. Thalaivar Biryani & Tiffin (Marathahalli) -> PENDING APPROVAL FOR ADMIN
		if err := seedResto(
			"Thalaivar Biryani & Tiffin - Marathahalli",
			"Tamil street food counter famous for hot iron-griddle egg kothu parotta and ghee podi idlis.",
			"Outer Ring Rd, Marathahalli Village, Bengaluru, Karnataka 560037",
			12.9535, 77.7012,
			"+91 80 4965 3322",
			"11:30 AM - 11:00 PM",
			restaurant.StatusPending,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Egg Kothu Parotta", "Shredded layered parotta tossed on hot iron griddle with egg and salna", "South Indian", 95.0, true},
				{"Hot Ghee Podi Idli (4 pcs)", "Tossed in fiery gun powder and clarified butter", "Breakfast", 50.0, true},
			},
		); err != nil {
			return err
		}

		// EDGE CASE 1: Meghana Foods (All items > ₹100, should be strictly EXCLUDED)
		if err := seedResto(
			"Meghana Foods - Residency Road",
			"Famous biryani chain known for spicy Andhra biryanis. All main items exceed ₹100.",
			"52, 1st Floor, Residency Rd, Bengaluru, Karnataka 560025",
			12.9719, 77.6068,
			"+91 80 4113 4567",
			"11:30 AM - 11:00 PM",
			restaurant.StatusApproved,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Meghana Special Chicken Biryani", "Boneless chicken marinated in spicy Andhra masala", "Biryani", 320.0, true},
				{"Paneer Biryani", "Rich aromatic basmati rice with spiced paneer cubes", "Biryani", 280.0, true},
				{"Chicken 65", "Spicy deep fried chicken chunks", "Starters", 310.0, true},
			},
		); err != nil {
			return err
		}

		// EDGE CASE 2: The Bangalore Cafe (Pending approval - should NOT appear to normal users)
		if err := seedResto(
			"The Bangalore Cafe",
			"Modern vegetarian fusion restaurant with global small plates.",
			"KH Rd, Shanthala Nagar, Bengaluru, Karnataka 560027",
			12.9644, 77.5991,
			"+91 80 2222 3456",
			"11:00 AM - 11:00 PM",
			restaurant.StatusPending,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Classic Fries", "Salted crispy potato fries", "Starters", 90.0, true},
				{"Cold Brew Coffee", "Steeped cold brew with tonic", "Beverages", 85.0, true},
				{"Mini Pizza Bites", "Bite sized cheesy snacks", "Snacks", 95.0, true},
			},
		); err != nil {
			return err
		}

		// EDGE CASE 3: Spice Terrace (Rejected listing - should NOT appear to normal users)
		if err := seedResto(
			"Spice Terrace",
			"Luxury poolside rooftop dining with North Indian kebabs.",
			"24, JW Marriott, Vittal Mallya Rd, Ashok Nagar, Bengaluru, Karnataka 560001",
			12.9719, 77.5956,
			"+91 80 6718 9999",
			"7:00 PM - 11:30 PM",
			restaurant.StatusRejected,
			[]struct {
				name, desc, category string
				price                float64
				avail                bool
			}{
				{"Masala Roasted Papad", "Crisp papad with spiced onion tomato salsa", "Starters", 90.0, true},
				{"Mint Chaas", "Spiced refreshing buttermilk", "Beverages", 80.0, true},
			},
		); err != nil {
			return err
		}

		// Link merchant owner to CTR restaurant if not yet assigned
		ctrRest, err := client.Restaurant.Query().Where(restaurant.NameEQ("Central Tiffin Room (CTR) - Shri Sagar")).Only(ctx)
		if err == nil && ctrRest != nil && merchantUser != nil {
			hasOwner, _ := ctrRest.QueryOwner().Exist(ctx)
			if !hasOwner {
				_, _ = ctrRest.Update().SetOwner(merchantUser).Save(ctx)
			}
		}

		return nil
}

