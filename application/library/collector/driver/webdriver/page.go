/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package webdriver

import "github.com/tebeka/selenium"

type Page struct {
	Driver selenium.WebDriver
}

func (s *Page) driver() selenium.WebDriver {
	return s.Driver
}

func (s *Page) FindElementById(locator string) selenium.WebElement {
	element, _ := s.Driver.FindElement(selenium.ByID, locator)
	return element
}

func (s *Page) FindElementByXpath(locator string) selenium.WebElement {
	element, _ := s.Driver.FindElement(selenium.ByXPATH, locator)
	return element
}

func (s *Page) FindElementByLinkText(locator string) selenium.WebElement {
	element, _ := s.Driver.FindElement(selenium.ByLinkText, locator)
	return element
}

func (s *Page) FindElementByPartialLink(locator string) selenium.WebElement {
	element, _ := s.Driver.FindElement(selenium.ByPartialLinkText, locator)
	return element
}

func (s *Page) FindElementByName(locator string) selenium.WebElement {
	element, _ := s.Driver.FindElement(selenium.ByName, locator)
	return element
}

func (s *Page) FindElementByTag(locator string) selenium.WebElement {
	element, _ := s.Driver.FindElement(selenium.ByTagName, locator)
	return element
}

func (s *Page) FindElementByClass(locator string) selenium.WebElement {
	element, _ := s.Driver.FindElement(selenium.ByClassName, locator)
	return element
}

func (s *Page) FindElementByCss(locator string) selenium.WebElement {
	element, _ := s.Driver.FindElement(selenium.ByCSSSelector, locator)
	return element
}

func (s *Page) FindElementsById(locator string) []selenium.WebElement {
	element, _ := s.Driver.FindElements(selenium.ByID, locator)
	return element
}

func (s *Page) FindElementsByXpath(locator string) []selenium.WebElement {
	element, _ := s.Driver.FindElements(selenium.ByXPATH, locator)
	return element
}

func (s *Page) FindElementsByLinkText(locator string) []selenium.WebElement {
	element, _ := s.Driver.FindElements(selenium.ByLinkText, locator)
	return element
}

func (s *Page) FindElementsByPartialLink(locator string) []selenium.WebElement {
	element, _ := s.Driver.FindElements(selenium.ByPartialLinkText, locator)
	return element
}

func (s *Page) FindElementsByName(locator string) []selenium.WebElement {
	element, _ := s.Driver.FindElements(selenium.ByName, locator)
	return element
}

func (s *Page) FindElementsByTag(locator string) []selenium.WebElement {
	element, _ := s.Driver.FindElements(selenium.ByTagName, locator)
	return element
}

func (s *Page) FindElementsByClass(locator string) []selenium.WebElement {
	element, _ := s.Driver.FindElements(selenium.ByClassName, locator)
	return element
}

func (s *Page) FindElementsByCss(locator string) []selenium.WebElement {
	element, _ := s.Driver.FindElements(selenium.ByCSSSelector, locator)
	return element
}

func (s *Page) MouseHoverToElement(locator string) selenium.WebElement {
	element, _ := s.Driver.FindElement(selenium.ByCSSSelector, locator)
	element.MoveTo(0, 0)
	return element
}
