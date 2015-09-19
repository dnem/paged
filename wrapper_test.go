package paged_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	. "github.com/dnem/paged"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResponseWrapper", func() {
	Context("when the message is wrapped in ErrorWrapper", func() {
		mx := mux.NewRouter()
		mx.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			err := errors.New("must feed a hamburger to the gnome before continuing")
			Formatter().JSON(w, http.StatusOK, ErrorWrapper(err.Error()))
			return
		}).Methods("GET")
		server := httptest.NewServer(mx)
		defer server.Close()

		url := server.URL + "/"
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		payload, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		var rw ResponseWrapper
		err = json.Unmarshal(payload, &rw)
		if err != nil {
			log.Fatal(err)
		}

		It("should have a default values", func() {
			Expect(rw.Status).To(Equal("error"))
			Expect(rw.Message).To(Equal("must feed a hamburger to the gnome before continuing"))
		})
	})

	Context("when the is wrapped with SuccessWrapper", func() {
		mx := mux.NewRouter()
		mx.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			params := ExtractRequestParams(req.URL.Query())
			Formatter().JSON(w, http.StatusOK, SuccessWrapper(&params))
			return
		}).Methods("GET")
		server := httptest.NewServer(mx)
		defer server.Close()

		url := server.URL + "/"
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		payload, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		var rw ResponseWrapper
		err = json.Unmarshal(payload, &rw)
		if err != nil {
			log.Fatal(err)
		}

		It("should have a default values", func() {
			Expect(rw.Status).To(Equal("success"))
			Expect(rw.Data).NotTo(BeNil())
			Expect(rw.Count).To(Equal(0))
			Expect(rw.Message).To(Equal(""))
		})
	})

	Context("when the is wrapped with CollectionWrapper", func() {
		mx := mux.NewRouter()
		mx.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			params := ExtractRequestParams(req.URL.Query())
			data := []int{1, 2, 3, 4, 5}
			count := len(data)
			Formatter().JSON(w, http.StatusOK, CollectionWrapper(data, count, params))
			return
		}).Methods("GET")
		server := httptest.NewServer(mx)
		defer server.Close()

		url := server.URL + "/?limit=2&offset=1&scope=fluffy&a=1&b=2"
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		payload, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		var rw ResponseWrapper
		err = json.Unmarshal(payload, &rw)
		if err != nil {
			log.Fatal(err)
		}

		It("should have a default values", func() {
			Expect(rw.Status).To(Equal("success"))
			Expect(rw.Data).NotTo(BeNil())
			Expect(rw.Count).To(Equal(5))
			Expect(rw.Message).To(Equal(""))
		})
	})
})
