package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
)

var filename, sheetname string

type details struct {
	Aff, Name, State, District, Principal, Phone, Email, Website string
}

var Data []details

func main() {

	fmt.Println("Enter the File Name :")
	fmt.Scan(&filename)
	fmt.Println("Enter the Sheet Name :")
	fmt.Scan(&sheetname)
	// Instantiate a new Colly collector
	c := colly.NewCollector()

	// Set the URL to send the form value request to
	url := "https://saras.cbse.gov.in/saras/AffiliatedList/ListOfSchdirReport"

	// Define the form data to be sent
	formData := map[string]string{
		"MainRadioValue": "Region_wise",
		"region":         "Trivendram",
	}

	// Attach a callback function to the OnResponse event
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status Code:", r.StatusCode)
	})

	c.OnHTML("tr", func(h *colly.HTMLElement) {
		var aff, state, district, name, principal, phone, email, website string
		h.ForEach("td", func(i int, c1 *colly.HTMLElement) {
			switch i {
			case 1:
				aff = c1.Text
				fmt.Println("affliation:", aff)
			case 2:
				place := c1.Text
				statePrefixPos := strings.Index(place, "State :")
				districtPrefixPos := strings.Index(place, "District :")

				// Extract the values for "KERALA" and "THRISSUR"
				state = strings.TrimSpace(place[statePrefixPos+len("State :") : districtPrefixPos])
				district = strings.TrimSpace(place[districtPrefixPos+len("District :"):])

				fmt.Println("state:", state, "\ndistrict:", district)
			case 4:
				name := c1.Text
				// Find the positions of "Name :" and "Head/Principal Name:"
				namePrefixPos := strings.Index(name, "Name :")
				headPrefixPos := strings.Index(name, "Head/Principal Name:")

				// Extract the values for "ST.LIOBA ACADEMY" and "JAIN MATHEW"
				schoolName := strings.TrimSpace(name[namePrefixPos+len("Name :") : headPrefixPos])
				principalName := strings.TrimSpace(name[headPrefixPos+len("Head/Principal Name:"):])

				name = schoolName
				principal = principalName
				fmt.Println("School Name:", name, "\nPrincipal Name:", principal)
				fmt.Println()
				fmt.Println()
			case 5:
				contact := c1.Text
				// Find the positions of "Phone No :", "Email :", and "Website :"
				phonePrefixPos := strings.Index(contact, "Phone No :")
				emailPrefixPos := strings.Index(contact, "Email :")
				websitePrefixPos := strings.Index(contact, "Website :")

				// Extract the values for the phone number, email, and website
				phone = strings.TrimSpace(contact[phonePrefixPos+len("Phone No :") : emailPrefixPos])
				email = strings.TrimSpace(contact[emailPrefixPos+len("Email :") : websitePrefixPos])
				website = strings.TrimSpace(contact[websitePrefixPos+len("Website :"):])

				fmt.Println("phone:", phone)
				fmt.Println("email:", email)
				fmt.Println("website:", website)
			default:
				fmt.Print()
			}
		})
		data := details{
			Aff:      aff,
			Name:     name,
			District: district,
			State:    state,
			Phone:    phone,
			Email:    email,
			Website:  website,
		}
		Data = append(Data, data)
	})

	// Send the form value request using the Post method
	c.Post(url, formData)
	writeXLSX(Data)
}

func writeXLSX(data []details) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(sheetname)
	if err != nil {
		panic(err)
	}

	row := sheet.AddRow()
	row.AddCell().SetValue("Affiliation No")
	row.AddCell().SetValue("Name")
	row.AddCell().SetValue("District")
	row.AddCell().SetValue("State")
	row.AddCell().SetValue("Phone")
	row.AddCell().SetValue("Email")
	row.AddCell().SetValue("Website")

	for _, r := range data {
		row := sheet.AddRow()
		row.AddCell().SetValue(r.Aff)
		row.AddCell().SetValue(r.Name)
		row.AddCell().SetValue(r.District)
		row.AddCell().SetValue(r.State)
		row.AddCell().SetValue(r.Phone)
		row.AddCell().SetValue(r.Email)
		row.AddCell().SetValue(r.Website)
	}

	err = file.Save(filename + ".xlsx")
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
