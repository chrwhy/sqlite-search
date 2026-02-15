package util

import (
	"math/rand"
	"sync"
	"time"
)

// ...existing code...

var rng *rand.Rand
var once sync.Once

func seed() {
	once.Do(func() {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
}

// RandInt returns a random integer in [min, max]. If min > max they are swapped.
func RandInt(min, max int) int {
	seed()
	if min > max {
		min, max = max, min
	}
	if min == max {
		return min
	}
	return rng.Intn(max-min+1) + min
}

// RandFloat64 returns a random float64 in [min, max).
func RandFloat64(min, max float64) float64 {
	seed()
	if min > max {
		min, max = max, min
	}
	return rng.Float64()*(max-min) + min
}

var firstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda", "William", "Elizabeth",
	"David", "Barbara", "Richard", "Susan", "Joseph", "Jessica", "Thomas", "Sarah", "Charles", "Karen",
	"Christopher", "Nancy", "Daniel", "Lisa", "Matthew", "Betty", "Anthony", "Margaret", "Donald", "Sandra",
	"Mark", "Ashley", "Paul", "Kimberly", "Steven", "Emily", "Andrew", "Donna", "Kenneth", "Michelle",
	"Joshua", "Carol", "George", "Amanda", "Kevin", "Melissa", "Brian", "Deborah", "Edward", "Stephanie",
	"Ronald", "Rebecca", "Timothy", "Laura", "Jason", "Sharon", "Jeffrey", "Cynthia", "Ryan", "Kathleen",
	"Jacob", "Amy", "Gary", "Shirley", "Nicholas", "Angela", "Eric", "Helen", "Stephen", "Anna",
	"Jonathan", "Brenda", "Larry", "Pamela", "Justin", "Nicole", "Scott", "Emma", "Brandon", "Samantha",
	"Benjamin", "Katherine", "Samuel", "Christine", "Frank", "Debra", "Gregory", "Rachel", "Raymond", "Carolyn",
	"Alexander", "Janet", "Patrick", "Catherine", "Jack", "Maria", "Dennis", "Heather", "Jerry", "Diane",
	"Tyler", "Julie", "Aaron", "Joyce", "Jose", "Victoria", "Adam", "Kelly", "Nathan", "Christina",
	"Zachary", "Lauren", "Harold", "Joan", "Christian", "Evelyn", "Keith", "Judith", "Roger", "Megan",
	"Gerald", "Andrea", "Ethan", "Cheryl", "Arthur", "Hannah", "Terry", "Jacqueline", "Lawrence", "Gloria",
	"Sean", "Ann", "Austin", "Teresa", "Carl", "Sara", "Albert", "Janice", "Dylan", "Jean",
	"Leo", "Alice", "Jordan", "Kathryn", "Roy", "Peggy", "Noah", "Sophie", "Billy", "Julia",
	"Bruce", "Ruby", "Willie", "Lois", "Alan", "Tina", "Eugene", "Nora", "Juan", "Ruth",
	"Wayne", "Marie", "Billy", "Eleanor", "Steve", "Eileen", "Louis", "Gail", "Jeremy", "Terence",
	"Fred", "Evelyn", "Philip", "Marilyn", "Bobby", "Opal", "Randy", "Lillian", "Howard", "Theresa",
	"Ethan", "June", "Vincent", "Doris", "Russell", "Belinda", "Louis", "Lori", "Phillip", "Wanda",
	"Earl", "Yvonne", "Craig", "Stacy", "Nathaniel", "Sally", "Caleb", "Mildred", "Cameron", "Beatrice",
	"Miguel", "Lydia", "Rafael", "Cecilia", "Calvin", "Christy", "Derek", "Harriet", "Jesse", "Tammy",
	"Oscar", "Marsha", "Victor", "Carmen", "Martin", "Gwendolyn", "Dustin", "Kara", "Travis", "Lorraine",
	"Curtis", "Kristin", "Clifford", "Darla", "Neil", "Nadine", "Mitchell", "Kelsey", "Gordon", "Priscilla",
	"Francis", "Monica", "Darren", "Holly", "Randall", "Tanya", "Dale", "Renee", "Shaun", "Martha",
	"Shane", "Lorena", "Dean", "Glenda", "Allan", "Adrienne", "Jon", "Sylvia", "Ruben", "Kristen",
	"Adrian", "Nina", "Clyde", "Brittany", "Glen", "Angelica", "Garry", "Belinda", "Derrick", "Clara",
	"Dwayne", "Yolanda", "Lance", "Celia", "Darren", "April", "Herman", "Marian", "Ramon", "Daisy",
	"Gilbert", "May", "Arnold", "Rita", "Ross", "Belle", "Edwin", "Adelaide", "Garry", "Cassandra",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
	"Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
	"Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson",
	"Walker", "Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores",
	"Green", "AdAMS", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell", "Carter", "Roberts",
	"Gomez", "Phillips", "Evans", "Turner", "Diaz", "Parker", "Cruz", "Edwards", "Collins", "Reyes",
	"Stewart", "Morris", "Morales", "Murphy", "Cook", "Rogers", "Gutierrez", "Ortiz", "Morgan", "Cooper",
	"Peterson", "Bailey", "Reed", "Kelly", "Howard", "Ramos", "Kim", "Cox", "Ward", "Richardson",
	"Watson", "Brooks", "Chavez", "Wood", "James", "Bennett", "Gray", "Mendoza", "Ruiz", "Hughes",
	"Price", "Alvarez", "Castillo", "Sanders", "Patel", "Myers", "Long", "Ross", "Foster", "Jimenez",
	"Powell", "Jenkins", "Perry", "Russell", "Sullivan", "Bell", "Coleman", "Butler", "Henderson", "Barnes",
	"Gonzales", "Fisher", "Vasquez", "Simmons", "Romero", "Jordan", "Patterson", "Alexander", "Hamilton", "Graham",
	"Reynolds", "Griffin", "Wallace", "Moreno", "West", "Cole", "Hayes", "Bryant", "Herrera", "Gibson",
	"Ellis", "Tran", "Medina", "Aguilar", "Stevens", "Murray", "Ford", "Castro", "Marshall", "Owens",
	"Harrison", "Fernandez", "Mcdonald", "Woods", "Washington", "Kennedy", "Wells", "Vargas", "Henry", "Chen",
	"Freeman", "Webb", "Tucker", "Guzman", "Burns", "Crawford", "Olson", "Simpson", "Porter", "Hunter",
	"Gordon", "Mendez", "Silva", "Shaw", "Snyder", "Mason", "Dixon", "Munoz", "Hunt", "Hicks",
	"Holmes", "Palmer", "Wagner", "Black", "Robertson", "Boyd", "Rose", "Stone", "Salazar", "Fox",
	"Weaver", "Baldwin", "Burnett", "Rowe", "Banks", "Meyer", "Bishop", "Mccoy", "Howell", "Alvarado",
	"Vega", "Chen", "Frederick", "Dunn", "Reeves", "Hudson", "Hamilton", "Spencer", "Lamb", "Carter",
	"Barker", "Gaines", "Clayton", "Leonard", "Walsh", "Lowe", "Schmidt", "Schneider", "Marshall", "Riley",
	"Hansen", "Cole", "Jensen", "Holt", "Gill", "Farmer", "Hart", "Warren", "Diaz", "Pena",
	"Richards", "Fitzgerald", "Mccarthy", "Valdez", "Weber", "Bates", "Miles", "Horton", "Nixon", "Hardy",
	"Hill", "Odonnell", "Nunez", "Mills", "Blair", "Nolan", "Mora", "Casey", "Roth", "Allison",
	"Sheppard", "Pace", "Horn", "English", "Atkinson", "Houston", "Rowland", "Baldwin", "Vance", "Bolton",
}

// RandomFirstName returns a random common English first name.
func RandomFirstName() string {
	seed()
	return firstNames[rng.Intn(len(firstNames))]
}

// RandomLastName returns a random common English last name.
func RandomLastName() string {
	seed()
	return lastNames[rng.Intn(len(lastNames))]
}

// RandomName returns a random "First Last" English name.
func RandomName() string {
	return RandomFirstName() + " " + RandomLastName()
}
