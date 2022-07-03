package utils

import "fmt"

func main() {
	b := []byte(`{"img":"https://distribution-point.webstorage-4sigma.it/donatazanotti-1149/media/progetto/2022/myneighbors3.jpg","format":[{"format":"webp","size":[{"width":"100","height":"100"},{"width":"200","height":"200"},{"width":"300","height":"300"},{"width":"500","height":"500"}]},{"format":"jpg","size":[{"width":"100","height":"100"},{"width":"200","height":"200"},{"width":"300","height":"300"},{"width":"500","height":"500"}]},{"format":"png","size":[{"width":"100","height":"100"},{"width":"200","height":"200"},{"width":"300","height":"300"},{"width":"500","height":"500"}]},{"format":"avif","size":[{"width":"100","height":"100"},{"width":"200","height":"200"},{"width":"300","height":"300"},{"width":"500","height":"500"}]}]}`)
	fmt.Println(b)
	// var imgData ImageGeneration
	// err := json.Unmarshal(b, &imgData)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// imageConverter(imgData)

}
