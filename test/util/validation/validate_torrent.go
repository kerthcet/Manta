/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"context"
	"errors"
	"fmt"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"

	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "github.com/inftyai/manta/api/v1alpha1"
	"github.com/inftyai/manta/test/util"
)

func ValidateTorrentStatusEqualTo(ctx context.Context, k8sClient client.Client, torrent *api.Torrent, conditionType string, reason string, status metav1.ConditionStatus) {
	gomega.Eventually(func() error {
		if err := k8sClient.Get(ctx, types.NamespacedName{Name: torrent.Name}, torrent); err != nil {
			return errors.New("failed to get torrent")
		}

		if torrent.Status.Phase == nil {
			return fmt.Errorf("status.phase should not be nil")
		}
		if *torrent.Status.Phase != conditionType {
			return fmt.Errorf("unexpected status.phase, want %s, got %s", conditionType, *torrent.Status.Phase)
		}

		if condition := apimeta.FindStatusCondition(torrent.Status.Conditions, conditionType); condition == nil {
			return fmt.Errorf("condition not found: %s", format.Object(torrent, 1))
		} else {
			if condition.Reason != reason || condition.Status != status {
				return fmt.Errorf("expected reason %q or status %q, but got %s", reason, status, format.Object(condition, 1))
			}
		}

		if torrent.Spec.ModelHub != nil && torrent.Spec.ModelHub.Filename != nil {
			if torrent.Status.Repo == nil || len(torrent.Status.Repo.Objects) != 1 {
				return fmt.Errorf("unexpected object length, should be equal to 1")
			}
		}

		if torrent.Spec.ModelHub != nil && torrent.Spec.ModelHub.Filename == nil {
			if torrent.Status.Repo == nil || len(torrent.Status.Repo.Objects) <= 1 {
				return fmt.Errorf("unexpected file length, should be greater than 1")
			}
		}

		return nil
	}, util.Timeout, util.Interval).Should(gomega.Succeed())
}

func ValidateAllReplicationsNodeNameEqualTo(ctx context.Context, k8sClient client.Client, torrent *api.Torrent, nodeName string) {
	gomega.Eventually(func() error {
		replicationList := api.ReplicationList{}
		selector := labels.SelectorFromSet(labels.Set{api.TorrentNameLabelKey: torrent.Name})
		if err := k8sClient.List(ctx, &replicationList, &client.ListOptions{
			LabelSelector: selector,
		}); err != nil {
			return err
		}

		for _, replication := range replicationList.Items {
			if replication.Name != nodeName {
				return fmt.Errorf("unexpected nodeName, expected %s, got %s", nodeName, replication.Name)
			}
		}
		return nil
	}, util.Timeout, util.Interval).Should(gomega.Succeed())
}

func ValidateReplicationsNumberEqualTo(ctx context.Context, k8sClient client.Client, torrent *api.Torrent, number int) {
	gomega.Eventually(func() bool {
		replicationList := api.ReplicationList{}
		selector := labels.SelectorFromSet(labels.Set{api.TorrentNameLabelKey: torrent.Name})
		if err := k8sClient.List(ctx, &replicationList, &client.ListOptions{
			LabelSelector: selector,
		}); err != nil {
			return false
		}

		return len(replicationList.Items) == number

	}, util.Timeout, util.Interval).Should(gomega.BeTrue())
}
