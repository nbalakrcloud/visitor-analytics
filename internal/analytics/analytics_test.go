package analytics_test

import (
	"testing"

	"github.com/nbalakrcloud/visitor-analytics/internal/analytics"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRecordVisitor(t *testing.T) {
	visitorStore := analytics.NewVisitorStore()
	visitorStore.RecordVisitor("abc.com", "1")
	Convey("With one visitor visiting one URL", t, func() {
		So(len(visitorStore.Store), ShouldEqual, 1)
		Convey("With two visitor visiting one URL", func() {
			visitorStore.RecordVisitor("abc.com", "2")
			So(len(visitorStore.Store), ShouldEqual, 1)
			Convey("With one URL visited by two visitors and one by only one visitor", func() {
				visitorStore.RecordVisitor("xyz.com", "1")
				So(len(visitorStore.Store), ShouldEqual, 2)
				So(len(visitorStore.Store["abc.com"]), ShouldEqual, 2)
				So(len(visitorStore.Store["xyz.com"]), ShouldEqual, 1)
			})
		})

	})

}

func TestGetUniqueVisitors(t *testing.T) {
	visitorStore := analytics.NewVisitorStore()
	visitorStore.RecordVisitor("abc.com", "1")
	Convey("With one visitor visiting one URL", t, func() {
		So(len(visitorStore.GetUniqueVisitors()), ShouldEqual, 1)
		Convey("With two visitor visiting one URL", func() {
			visitorStore.RecordVisitor("abc.com", "2")
			So(len(visitorStore.GetUniqueVisitors()), ShouldEqual, 1)
			Convey("With one URL visited by two visitors and one by only one visitor", func() {
				visitorStore.RecordVisitor("xyz.com", "1")
				So(len(visitorStore.GetUniqueVisitors()), ShouldEqual, 2)
				So(visitorStore.GetUniqueVisitors()["abc.com"], ShouldEqual, 2)
				So(visitorStore.GetUniqueVisitors()["xyz.com"], ShouldEqual, 1)
			})
		})

	})

}
